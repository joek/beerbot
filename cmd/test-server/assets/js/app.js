



$(function () {
  console.log("touchscreen is", VirtualJoystick.touchScreenAvailable() ? "available" : "not available");
  var conn;
  var joystick	= new VirtualJoystick({
    mouseSupport: true,
    stationaryBase: true,
    baseX: 200,
    baseY: 200,
    limitStickTravel: true,
    stickRadius: 100
  });

  joystick.addEventListener('touchStart', function(){
    console.log('down')
    connectRobot()
  })
  joystick.addEventListener('touchEnd', function(){
    console.log('up')
    disconnectRobot()
  })

  joystick.addEventListener('')

  $( "body" )
  .mouseup(function() {
    console.log('up')
    disconnectRobot()
  })
  .mousedown(function() {
    console.log('down')
    connectRobot()
  });


  setInterval(function(){
    x = joystick.deltaX();
    y = joystick.deltaY();

    t = calculateTank(x, y);
    sendCommand({motor: t})
  }, 1/30 * 1000);

  function connectRobot(){
    console.log("Connect websocket")
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + window.location.host + "/ws");
        conn.onclose = function (evt) {
            console.log("Connection closed.");
        };
    } else {
      alert("Sorry, your browser doesn't support Device Orientation");
    }
  }

  function disconnectRobot(){
    console.log("Disconnect websocket")
    if(conn){
      conn.close()
      conn = null
    }
  }

  function sendCommand(command) {
    var message = JSON.stringify(command);
    if (!conn) {
      return;
    }
    console.log("Send: " + message)
    conn.send(message);
  }

  function calculateTank(x, y){
    // First hypotenuse
    var z = Math.sqrt(x*x + y*y);
    // angle in radians
    rad = Math.acos(Math.abs(x)/z);
    // and in degrees
    angle = rad*180/Math.PI;
    // Now angle indicates the measure of turn
    // Along a straight line, with an angle o, the turn co-efficient is same
    // this applies for angles between 0-90, with angle 0 the co-eff is -1
    // with angle 45, the co-efficient is 0 and with angle 90, it is 1
    var tcoeff = -1 + (angle/90)*2;
    var turn = tcoeff * Math.abs(Math.abs(y) - Math.abs(x));
    turn = Math.round(turn*100)/100;
    // And max of y or x is the movement
    var move = Math.max(Math.abs(y),Math.abs(x));

    // First and third quadrant
    if( (x >= 0 && y >= 0) || (x < 0 &&  y < 0) ) {
        right = move;
        left = turn;
    } else {
        left = move;
        right = turn;
    }

    // Reverse polarity
    if(y > 0) {
        left = 0 - left;
        right = 0 - right;
    }

    return { left: left / 100, right: right / 100}
  }
});
