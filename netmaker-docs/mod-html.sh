#!/bin/bash


sed -i '/rel=\"next\"\ title=\"About\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/\"\ />' /usr/share/nginx/html/index.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/server-installation.html\"\ />' /usr/share/nginx/html/server-installation.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/license.html\"\ />' /usr/share/nginx/html/license.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/pro/pro-metrics.html\"\ />' /usr/share/nginx/html/pro/pro-metrics.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/pro/pro-setup.html\"\ />' /usr/share/nginx/html/pro/pro-setup.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/pro/pro-users.html\"\ />' /usr/share/nginx/html/pro/pro-users.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/pro/pro-branding.html\"\ />' /usr/share/nginx/html/pro/pro-branding.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/pro/index.html\"\ />' /usr/share/nginx/html/pro/index.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/pro/rac.html\"\ />' /usr/share/nginx/html/pro/rac.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/about.html\"\ />' /usr/share/nginx/html/about.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/acls.html\"\ />' /usr/share/nginx/html/acls.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/api.html\"\ />' /usr/share/nginx/html/api.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/conduct.html\"\ />' /usr/share/nginx/html/conduct.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/architecture.html\"\ />' /usr/share/nginx/html/architecture.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/troubleshoot.html\"\ />' /usr/share/nginx/html/troubleshoot.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/nmctl.html\"\ />' /usr/share/nginx/html/nmctl.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/nmctl-standalone.html\"\ />' /usr/share/nginx/html/nmctl-standalone.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/usage.html\"\ />' /usr/share/nginx/html/usage.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/external-clients.html\"\ />' /usr/share/nginx/html/external-clients.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/egress-gateway.html\"\ />' /usr/share/nginx/html/egress-gateway.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/netclient.html\"\ />' /usr/share/nginx/html/netclient.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/quick-start.html\"\ />' /usr/share/nginx/html/quick-start.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/support.html\"\ />' /usr/share/nginx/html/support.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/upgrades.html\"\ />' /usr/share/nginx/html/upgrades.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/getting-started.html\"\ />' /usr/share/nginx/html/getting-started.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/advanced-client-install.html\"\ />' /usr/share/nginx/html/advanced-client-install.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/install.html\"\ />' /usr/share/nginx/html/install.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/oauth.html\"\ />' /usr/share/nginx/html/oauth.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/pro/pro-relay-server.html\"\ />' /usr/share/nginx/html/pro/pro-relay-server.html

sed -i '/link rel=\"prev\"/a \ \ \ \ <link rel=\"canonical\"\ href=\"https://docs.netmaker.io/ui-reference.html\"\ />' /usr/share/nginx/html/ui-reference.html

pushd /usr/share/nginx/html/
find . -type f -name "*.html" -exec sed -i 's-var json_loc = "../"versions.json""-var json_loc = "../versions.json"-' {} +
find . -type f -name "*.html" -exec sed -i 's-var json_loc = ""versions.json""-var json_loc = "versions.json"-' {} +
find . -type f -name "*.html" -exec sed -i 's-</head>-<link rel="shortcut icon" href="#">\n</head>-' {} +
find . -type f -name "*.html" -exec sed -i 's#</head>#<script async src="REPLACETHISTHING" referrerpolicy="strict-origin-when-cross-origin"></script>\n</head>#' {} +
find . -type f -name "*.html" -exec sed -i """s#</head>#<script> (function(ss,ex){ window.ldfdr=window.ldfdr||function(){(ldfdr._q=ldfdr._q||[]).push([].slice.call(arguments));}; (function(d,s){ fs=d.getElementsByTagName(s)[0]; function ce(src){ var cs=d.createElement(s); cs.src=src; cs.async=1; fs.parentNode.insertBefore(cs,fs); }; ce('https://sc.lfeeder.com/lftracker_v1_'+ss+(ex?'_'+ex:'')+'.js'); })(document,'script'); })('p1e024BeABJ7GB6d'); </script>\n</head>#""" {} +
find . -type f -name "*.html" -exec sed -i """s#</head>#<script>
    !function(t,e){var o,n,p,r;e.__SV||(window.posthog=e,e._i=[],e.init=function(i,s,a){function g(t,e){var o=e.split(\".\");2==o.length&&(t=t[o[0]],e=o[1]),t[e]=function(){t.push([e].concat(Array.prototype.slice.call(arguments,0)))}}(p=t.createElement(\"script\")).type=\"text/javascript\",p.async=!0,p.src=s.api_host.replace(\".i.posthog.com\",\"-assets.i.posthog.com\")+\"/static/array.js\",(r=t.getElementsByTagName(\"script\")[0]).parentNode.insertBefore(p,r);var u=e;for(void 0!==a?u=e[a]=[]:a=\"posthog\",u.people=u.people||[],u.toString=function(t){var e=\"posthog\";return\"posthog\"!==a&&(e+=\".\"+a),t||(e+=\" (stub)\"),e},u.people.toString=function(){return u.toString(1)+\".people (stub)\"},o=\"capture identify alias people.set people.set_once set_config register register_once unregister opt_out_capturing has_opted_out_capturing opt_in_capturing reset isFeatureEnabled onFeatureFlags getFeatureFlag getFeatureFlagPayload reloadFeatureFlags group updateEarlyAccessFeatureEnrollment getEarlyAccessFeatures getActiveMatchingSurveys getSurveys getNextSurveyStep onSessionId setPersonProperties\".split(\" \"),n=0;n<o.length;n++)g(u,o[n]);e._i.push([i,s,a])},e.__SV=1)}(document,window.posthog||[]);
    posthog.init('POSTHOG_KEY',{api_host:'https://us.i.posthog.com', person_profiles: 'identified_only'})
</script>\n</head>#""" {} +
popd
