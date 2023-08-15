/*
Copyright © 2023 Netmaker Team

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gravitl/netclient/config"
	"github.com/gravitl/netclient/ncutils"
	"github.com/gravitl/netmaker/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"

	_ "github.com/mattn/go-sqlite3" // need to blank import this package
)

var rootCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade server netclient config",
	Long:  `upgrade server netclient config`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		upgrade()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	logLevel := &slog.LevelVar{}
	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			a.Value = slog.StringValue(filepath.Base(a.Value.String()))
		}
		return a
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replace, Level: logLevel}))
	slog.SetDefault(logger)
}

func upgrade() {
	slog.Info("creating config dir")
	if err := os.MkdirAll("/etc/netclient/config", os.ModePerm); err != nil {
		slog.Error("create dir", "error", err)
	}
	slog.Info("retrieving legacy nodes")
	serverPassword := ncutils.RandomString(32)
	nodes, err := getAllLegacyNodes(serverPassword)
	if err != nil {
		slog.Error("unable to get nodes", "error", err)
		os.Exit(1)
	}
	for _, node := range nodes {
		if node.IsServer != "yes" {
			continue
		}
		slog.Info("saving password")
		if err := os.WriteFile("/etc/netclient/config/secret-"+node.Network, []byte(serverPassword), 0600); err != nil {
			slog.Error("saving password", "error", err)
		}
		slog.Info("saving traffic keys")
		node.TrafficKeys.Mine, err = setTrafficKeys(node.Network)
		if err != nil {
			slog.Error("traffickeys", err)
		}
		slog.Info("saving wg key")
		saveWGPrivateKey(node.ID, node.Network)
		serverTrafficKey, err := getServerTrafficKey()
		if err != nil {
			slog.Error("get server traffic key", "error", err)
		}
		slog.Info("saving node", "name", node.Name)
		saveNode(node, serverTrafficKey)
	}
}

func getAllLegacyNodes(password string) ([]models.LegacyNode, error) {
	var key, value string
	nodes := []models.LegacyNode{}
	node := models.LegacyNode{}
	db, err := sql.Open("sqlite3", "/var/lib/docker/volumes/root_sqldata/_data/netmaker.db")
	if err != nil {
		return nodes, fmt.Errorf("open database %w", err)
	}
	row, err := db.Query("SELECT * from nodes")
	if err != nil {
		return nodes, fmt.Errorf("db query %w", err)
	}
	records := make(map[string]string)
	defer row.Close()
	for row.Next() {
		row.Scan(&key, &value)
		records[key] = value
	}
	if len(records) == 0 {
		return nodes, errors.New("no records")
	}
	for key, data := range records {
		if err := json.Unmarshal([]byte(data), &node); err != nil {
			slog.Warn("unmarhal node", "error", err)
			continue
		}
		if node.IsServer != "yes" {
			continue
		}
		slog.Info("processing node", "name", node.Name)
		//update password
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 5)
		if err != nil {
			slog.Warn("bcrypt", "error", err)
			continue
		}
		node.Password = password
		nodes = append(nodes, node)
		node.Password = string(hash)
		nodeUpdate, err := json.Marshal(node)
		if err != nil {
			slog.Warn("marshal node", "error", err)
			continue
		}
		if _, err := db.Exec("INSERT or REPLACE INTO nodes (key, value) VALUES (?, ?)", key, string(nodeUpdate)); err != nil {
			slog.Warn("update node", "error", err)
			continue
		}
	}
	return nodes, nil

}

func setTrafficKeys(network string) ([]byte, error) {
	slog.Info("setting traffic keys")
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return []byte{}, fmt.Errorf("generate trafffickeys %w", err)
	}
	privBytes, err := ncutils.ConvertKeyToBytes(priv)
	if err != nil {
		return []byte{}, fmt.Errorf("convert private traffickey %w", err)
	}
	if err := os.WriteFile("/etc/netclient/config/traffic-"+network, privBytes, 0600); err != nil {
		return []byte{}, fmt.Errorf("saving traffickey %w", err)
	}
	pubBytes, err := ncutils.ConvertKeyToBytes(pub)
	if err != nil {
		return []byte{}, fmt.Errorf("convert pub traffickey %w", err)
	}
	return pubBytes, nil
}

func saveWGPrivateKey(id, network string) {
	type Key struct {
		PrivateKey string
	}
	var value string
	var key Key
	db, err := sql.Open("sqlite3", "/var/lib/docker/volumes/root_sqldata/_data/netmaker.db")
	if err != nil {
		return
	}
	row, err := db.Query("SELECT VALUE from serverconf where key='" + id + "'")
	if err != nil {
		return
	}
	row.Next()
	row.Scan(&value)
	db.Close()
	if err := json.Unmarshal([]byte(value), &key); err != nil {
		return
	}
	if err := os.WriteFile("/etc/netclient/config/wgkey-"+network, []byte(key.PrivateKey), 0600); err != nil {
		return
	}
	return
}

func getServerTrafficKey() ([]byte, error) {
	type Key struct {
		UUID           string `json:"uuid"`
		LastSend       int64  `json:"lastsend"`
		TrafficKeyPriv []byte `json:"traffickeypriv"`
		TrafficKeyPub  []byte `json:"traffickeypub"`
	}
	var value string
	var key Key
	db, err := sql.Open("sqlite3", "/var/lib/docker/volumes/root_sqldata/_data/netmaker.db")
	if err != nil {
		return []byte{}, err
	}
	row, err := db.Query("SELECT VALUE from serveruuid where key='serveruuid'")
	if err != nil {
		return []byte{}, err
	}
	row.Next()
	row.Scan(&value)
	db.Close()
	if err := json.Unmarshal([]byte(value), &key); err != nil {
		return []byte{}, err
	}
	return key.TrafficKeyPub, nil
}

func saveNode(node models.LegacyNode, serverTrafficKey []byte) error {
	var cfg config.ClientConfig
	cfg.Node = node
	cfg.Node.TrafficKeys.Server = serverTrafficKey
	cfg.NetworkSettings = node.NetworkSettings
	cfg.Network = node.Network
	f, err := os.OpenFile("/etc/netclient/config/netconfig-"+node.Network, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open node file %w", err)
	}
	defer f.Close()
	compose, err := os.ReadFile("./docker-compose.yml")
	if err != nil {
		return fmt.Errorf("read compose file %w", err)
	}
	reg := regexp.MustCompile(`.*SERVER_NAME:.*`)
	matches := reg.FindStringSubmatch(string(compose))
	if len(matches) == 0 {
		return errors.New("no server name")
	}
	server := strings.ReplaceAll(matches[0], "SERVER_NAME: ", "")
	server = strings.ReplaceAll(server, "\"", "")
	server = strings.ReplaceAll(server, " ", "")
	cfg.Server.Server = server
	if err := yaml.NewEncoder(f).Encode(cfg); err != nil {
		return fmt.Errorf("encode node file %w", err)
	}
	return f.Sync()

}
