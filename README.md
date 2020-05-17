# gelftest
Utility to test connections to graylog using the GELF TCP protocol

# usage:
```
Usage of gelftest:
gelftest [options] [message]
  -c int
        Number of messages to send (shorthand) (default 1)
  -count int
        Number of messages to send (default 1)
  -e string
        the source_env field in the message (shorthand) (default "dev")
  -g string
        The Graylog server (shorthand) (default "localhost")
  -graylog string
        The Graylog server (default "localhost")
  -logtype string
        The logtype (APP or EVENT) (default "APP")
  -p int
        The port of the Graylog server (shorthand) (default 12201)
  -port int
        The port of the Graylog server (default 12201)
  -s int
        Sleeptime in milliseconds between sends
  -sourceenv string
        the source_env field in the message (default "dev")
  -t string
        The logtype (APP or EVENT) (shorthand) (default "APP")
  -v    Be verbose
  -verbose
        Be verbose
```
