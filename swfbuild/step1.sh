#!/bin/bash

echo "STEP1: export xml"
if [ ! -f ja.swf ]; then
	curl http://inishie-dungeon.com/inishie.swf > ja.swf
fi

if [ ! -f en.swf ]; then
	curl http://inidunres.0w0.be/inishie.swf > en.swf
fi

FFDEC_HOME=D:/usr/soft/ffdec

java -jar $FFDEC_HOME/ffdec.jar -swf2xml ja.swf ja.xml
java -jar $FFDEC_HOME/ffdec.jar -swf2xml en.swf en.xml

# rm -rf sources
# mkdir sources
# java -jar $FFDEC_HOME/ffdec.jar -config autoDeobfuscate=1,parallelSpeedUp=0 -export script sources ja.swf
# java -jar $FFDEC_HOME/ffdec.jar -format fla:cs5.5 -export fla sources/myfile.fla ja.swf