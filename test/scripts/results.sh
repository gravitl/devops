#! /bin/bash
if [ -f /tmp/clean/log ]
then
    grep ERROR /tmp/clean/clean.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "removing gateways: FAILED" >> /tmp/results.log
    else
        echo "removing gateways: PASSED" >> /tmp/results.log
    fi
else
    echo "removing gateways: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/ping/ping.log ]
then
    grep ERROR /tmp/ping/ping.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "initial ping: FAILED" >> /tmp/results.log
    else
        echo "initial ping: PASSED" >> /tmp/results.log
    fi
else
    echo "initial ping: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/tests/peerupdate.log ]
then
    grep ERROR /tmp/tests/peerupdate.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "peerupdate: FAILED" >> /tmp/results.log
    else
        echo "peerupdate: PASSED" >> /tmp/results.log
    fi
else
    echo "peerupdate: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/tests/ingress.log ]
then
    grep ERROR /tmp/tests/ingress.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "ingress: FAILED" >> /tmp/results.log
    else
        echo "ingress: PASSED" >> /tmp/results.log
    fi
else
    echo "ingress: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/test/egress.log ]
then
    grep ERROR /tmp/tests/egress.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "egress: FAILED" >> /tmp/results.log
    else
        echo "egress: PASSED" >> /tmp/results.log
    fi
else
    echo "egress: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/tests/relay.log ]
then
    grep ERROR /tmp/tests/relay.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "relay: FAILED" >> /tmp/results.log
    else
        echo "relay: PASSED" >> /tmp/results.log
    fi
else
    echo "relay: NOT RUN" >> /tmp/results.log
fi

if [ -f /tmp/ping2/ping.log ]
then
    grep ERROR /tmp/ping2/ping.log >> /tmp/errors.log
    if [ $? -eq 0 ]
    then
        echo "final ping: FAILED" >> /tmp/results.log
    else
        echo "final ping: PASSED" >> /tmp/results.log
    fi
else
    echo "final ping: NOT RUN" >> /tmp/results.log
fi
