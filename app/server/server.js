var express = require("express");
var app = express();
app.use(express.static(__dirname + '/../dist/process-manager-client'));
var server = app.listen(8080, function(){
  console.log("Started. Port:" + server.address().port);
});
