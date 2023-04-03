// Import dependencies
const express = require("express");
const request = require("request");
const ejs = require("ejs");

// Create Express app
const app = express();

// Set up EJS as the view engine
app.set("view engine", "ejs");

// Set up static directory for CSS files
app.use(express.static(__dirname + "/public"));

// Define route for index page
app.get("/", function (req, res) {
  // Make API request to get data
  request(
    //AppConfig
    "http://localhost:2772/applications/blogAppConfigGo/environments/prod/configurations/whichSide?flag=allegiance",
    //Paladin
    // "https://mocki.io/v1/58155774-c261-4e1e-a2a9-576bf9e9977c",
    //Dark Knight
    // "https://mocki.io/v1/4c6a46fc-011e-4a3b-b621-174fd9778058",
    function (error, response, body) {
      if (!error && response.statusCode == 200) {
        // Parse JSON data
        const data = JSON.parse(body);
        // Render index page using EJS template
        res.render("index", { data: data });
      } else {
        // Handle error
        res.render("error");
      }
    }
  );
});

// Start server
const port = process.env.PORT || 8088;
app.listen(port, function () {
  console.log("Server running on port " + port);
});