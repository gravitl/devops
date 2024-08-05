#! /bin/bash

# clean results for 3 networks
if [ -f /tmp/clean-devops/clean.log ]
then
    grep ERROR /tmp/clean-devops/clean.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "removing gateways on devops: FAILED" >> /tmp/results.log
    else
        echo "removing gateways on devops: PASSED" >> /tmp/results.log
    fi
else
    echo "removing gateways on devops: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/clean-devops4/clean.log ]
then
    grep ERROR /tmp/clean-devops4/clean.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "removing gateways on devops4: FAILED" >> /tmp/results.log
    else
        echo "removing gateways on devops4: PASSED" >> /tmp/results.log
    fi
else
    echo "removing gateways on devops4: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/clean-devopsv6/clean.log ]
then
    grep ERROR /tmp/clean-devopsv6/clean.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "removing gateways on devopsv6: FAILED" >> /tmp/results.log
    else
        echo "removing gateways on devopsv6: PASSED" >> /tmp/results.log
    fi
else
    echo "removing gateways on devopsv6: NOT RUN" >> /tmp/results.log
fi

# ping results for 3 networks

if [ -f /tmp/ping-devops/ping.log ]
then
    grep ERROR /tmp/ping-devops/ping.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "initial ping on devops: FAILED" >> /tmp/results.log
    else
        echo "initial ping on devops: PASSED" >> /tmp/results.log
    fi
else
    echo "initial ping on devops: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/ping-devops4/ping.log ]
then
    grep ERROR /tmp/ping-devops4/ping.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "initial ping on devops4: FAILED" >> /tmp/results.log
    else
        echo "initial ping on devops4: PASSED" >> /tmp/results.log
    fi
else
    echo "initial ping on devops4: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/ping-devopsv6/ping.log ]
then
    grep ERROR /tmp/ping-devopsv6/ping.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "initial ping on devopsv6: FAILED" >> /tmp/results.log
    else
        echo "initial ping on devopsv6: PASSED" >> /tmp/results.log
    fi
else
    echo "initial ping on devopsv6: NOT RUN" >> /tmp/results.log
fi

# test results for 3 networks

if [ -f /tmp/tests-devops/peerupdate.log ]
then
    grep ERROR /tmp/tests-devops/peerupdate.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "peerupdate for devops: FAILED" >> /tmp/results.log
    else
        echo "peerupdate for devops: PASSED" >> /tmp/results.log
    fi
else
    echo "peerupdate for devops: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/tests-devops4/peerupdate.log ]
then
    grep ERROR /tmp/tests-devops4/peerupdate.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "peerupdate for devops4: FAILED" >> /tmp/results.log
    else
        echo "peerupdate for devops4: PASSED" >> /tmp/results.log
    fi
else
    echo "peerupdate for devops4: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/tests-devopsv6/peerupdate.log ]
then
    grep ERROR /tmp/tests-devopsv6/peerupdate.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "peerupdate for devopsv6: FAILED" >> /tmp/results.log
    else
        echo "peerupdate for devopsv6: PASSED" >> /tmp/results.log
    fi
else
    echo "peerupdate for devopsv6: NOT RUN" >> /tmp/results.log
fi

# ingress results for 3 networks

if [ -f /tmp/tests-devops/ingress.log ]
then
    grep ERROR /tmp/tests-devops/ingress.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "ingress for devops: FAILED" >> /tmp/results.log
    else
        echo "ingress for devops: PASSED" >> /tmp/results.log
    fi
else
    echo "ingress for devops: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/tests-devops4/ingress.log ]
then
    grep ERROR /tmp/tests-devops4/ingress.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "ingress for devops4: FAILED" >> /tmp/results.log
    else
        echo "ingress for devops4: PASSED" >> /tmp/results.log
    fi
else
    echo "ingress for devops4: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/tests-devopsv6/ingress.log ]
then
    grep ERROR /tmp/tests-devopsv6/ingress.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "ingress for devopsv6: FAILED" >> /tmp/results.log
    else
        echo "ingress for devopsv6: PASSED" >> /tmp/results.log
    fi
else
    echo "ingress for devopsv6: NOT RUN" >> /tmp/results.log
fi

# egress results

if [ -f /tmp/tests-devops/egress.log ]
then
    grep ERROR /tmp/tests-devops/egress.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "egress for devops: FAILED" >> /tmp/results.log
    else
        echo "egress for devops: PASSED" >> /tmp/results.log
    fi
else
    echo "egress for devops: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/tests-devops4/egress.log ]
then
    grep ERROR /tmp/tests-devops4/egress.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "egress for devops4: FAILED" >> /tmp/results.log
    else
        echo "egress for devops4: PASSED" >> /tmp/results.log
    fi
else
    echo "egress for devops4: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/tests-devopsv6/egress.log ]
then
    grep ERROR /tmp/tests-devopsv6/egress.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "egress for devopsv6: FAILED" >> /tmp/results.log
    else
        echo "egress for devopsv6: PASSED" >> /tmp/results.log
    fi
else
    echo "egress for devopsv6: NOT RUN" >> /tmp/results.log
fi



# relay results

if [ -f /tmp/tests-devops/relay.log ]
then
    grep ERROR /tmp/tests-devops/relay.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "relay for devops: FAILED" >> /tmp/results.log
    else
        grep "WARN" /tmp/tests-devops/relay.log
        if [ $? -eq 0 ]
        then
            echo "relay for devops: SKIPPED" >> /tmp/results.log
        else
            echo "relay for devops: PASSED" >> /tmp/results.log
        fi
    fi
else
    echo "relay for devops: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/tests-devops4/relay.log ]
then
    grep ERROR /tmp/tests-devops4/relay.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "relay for devops4: FAILED" >> /tmp/results.log
    else
        grep "WARN" /tmp/tests-devops4/relay.log
        if [ $? -eq 0 ]
        then
            echo "relay for devops4: SKIPPED" >> /tmp/results.log
        else
            echo "relay for devops4: PASSED" >> /tmp/results.log
        fi
    fi
else
    echo "relay for devops4: NOT RUN" >> /tmp/results.log
fi


if [ -f /tmp/tests-devopsv6/relay.log ]
then
    grep ERROR /tmp/tests-devopsv6/relay.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "relay for devopsv6: FAILED" >> /tmp/results.log
    else
        grep "WARN" /tmp/tests-devopsv6/relay.log
        if [ $? -eq 0 ]
        then
            echo "relay for devopsv6: SKIPPED" >> /tmp/results.log
        else
            echo "relay for devopsv6: PASSED" >> /tmp/results.log
        fi
    fi
else
    echo "relay for devopsv6: NOT RUN" >> /tmp/results.log
fi


# internet gateway results

if [ -f /tmp/tests-devops/internetGateway.log ]
then
    grep ERROR /tmp/tests-devops/internetGateway.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "internet gateway for devops: FAILED" >> /tmp/results.log
    else
        grep "WARN" /tmp/tests-devops/internetGateway.log
        if [ $? -eq 0 ]
        then
            echo "internet gateway for devops: SKIPPED" >> /tmp/results.log
        else
            echo "internet gateway for devops: PASSED" >> /tmp/results.log
        fi
    fi
else
    echo "internet gateway for devops: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/tests-devops4/internetGateway.log ]
then
    grep ERROR /tmp/tests-devops4/internetGateway.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "internet gateway for devops4: FAILED" >> /tmp/results.log
    else
        grep "WARN" /tmp/tests-devops4/internetGateway.log
        if [ $? -eq 0 ]
        then
            echo "internet gateway for devops4: SKIPPED" >> /tmp/results.log
        else
            echo "internet gateway for devops4: PASSED" >> /tmp/results.log
        fi
    fi
else
    echo "internet gateway for devops4: NOT RUN" >> /tmp/results.log
fi


if [ -f /tmp/tests-devopsv6/internetGateway.log ]
then
    grep ERROR /tmp/tests-devopsv6/internetGateway.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "internet gateway for devopsv6: FAILED" >> /tmp/results.log
    else
        grep "WARN" /tmp/tests-devopsv6/internetGateway.log
        if [ $? -eq 0 ]
        then
            echo "internet gateway for devopsv6: SKIPPED" >> /tmp/results.log
        else
            echo "internet gateway for devopsv6: PASSED" >> /tmp/results.log
        fi
    fi
else
    echo "internet gateway for devopsv6: NOT RUN" >> /tmp/results.log
fi




# final ping results

if [ -f /tmp/ping2-devops/ping.log ]
then
    grep ERROR /tmp/ping2-devops/ping.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "final ping for devops: FAILED" >> /tmp/results.log
    else
        echo "final ping for devops: PASSED" >> /tmp/results.log
    fi
else
    echo "final ping for devops: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/ping2-devops4/ping.log ]
then
    grep ERROR /tmp/ping2-devops4/ping.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "final ping for devops4: FAILED" >> /tmp/results.log
    else
        echo "final ping for devops4: PASSED" >> /tmp/results.log
    fi
else
    echo "final ping for devops4: NOT RUN" >> /tmp/results.log
fi


if [ -f /tmp/ping2-devopsv6/ping.log ]
then
    grep ERROR /tmp/ping2-devopsv6/ping.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "final ping for devopsv6: FAILED" >> /tmp/results.log
    else
        echo "final ping for devopsv6: PASSED" >> /tmp/results.log
    fi
else
    echo "final ping for devopsv6: NOT RUN" >> /tmp/results.log
fi

