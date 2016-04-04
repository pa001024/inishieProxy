#!/bin/bash

echo "viewcode"
if [ ! -f ja.swf ]; then
	curl http://inishie-dungeon.com/inishie.swf > ja.swf
fi

FFDEC_HOME=D:/usr/soft/ffdec

cp en.swf en-kai.swf
java -jar $FFDEC_HOME/ffdec.jar en-kai.swf

# rm -rf sources
# mkdir sources
# java -jar $FFDEC_HOME/ffdec.jar -config autoDeobfuscate=1,parallelSpeedUp=0 -export script sources ja.swf
# java -jar $FFDEC_HOME/ffdec.jar -format fla:cs5.5 -export fla sources/myfile.fla ja.swf