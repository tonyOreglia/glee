## to stream engine messages to terminal
cat < engine.log
## to send the engine a message from terminal
echo isready > ui.log

## After setting up pipes via ... 	
wPipe := pipe.NewPipe("engine.log")
rPipe := pipe.NewPipe("ui.log")