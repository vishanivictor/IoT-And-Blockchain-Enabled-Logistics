#include <ESP8266HTTPClient.h>
#include <ESP8266WiFi.h>

int outputpin= A0;// initializes the output pin of the LM35 temperature sensor
#define URL "http://<<ip>>:8080/sensorreading"
 
void setup() {
 
  Serial.begin(115200);//Serial connection
  WiFi.begin("<<user-name>>", "<<password>>"); //WiFi connection
 
  while (WiFi.status() != WL_CONNECTED) {  //Wait for the WiFI connection completion
 
    delay(500);
    Serial.println("Waiting for connection");
 
  }
 
}
 
void loop() {
 
 if(WiFi.status()== WL_CONNECTED){   //Check WiFi connection status

    int analogTemp = analogRead(outputpin);
    float millivolts = (analogTemp/1024.0) * 3300; //3300 is the voltage provided by NodeMCU
    float celsius = millivolts/10;
    Serial.println(celsius);
    
    String payload = "{\"sensor\":\"sensorA\"";
    payload += ",\"temp\":\"";
    payload += celsius;
    payload += "\"}";
 
    Serial.print("Sending payload: ");
    Serial.println(payload);
    
    HTTPClient http;    //Declare object of class HTTPClient
 
    http.begin(URL);      //Specify request destination
    http.addHeader("Content-Type", "application/json");  //Specify content-type header
 
    int httpCode = http.POST(payload);   //Send the request
    String response = http.getString();                  //Get the response payload
 
    Serial.println(httpCode);   //Print HTTP return code
    Serial.println(response);    //Print request response payload
 
    http.end();  //Close connection
 
 }else{
 
    Serial.println("Error in WiFi connection");   
 
 }

 //Send temperature reading at every 5 minutes
 for(int i=0;i<60*5;i++)
   delay(1000);  
 
}
