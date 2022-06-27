const http = require('http');
const AWS = require('aws-sdk');
const docClient = new AWS.DynamoDB.DocumentClient();

exports.handler = async (event) => {

    const res = await new Promise((resolve, reject) => {
        http.get(
            "http://localhost:2772/applications/AWSomeCribRentals/environments/Beta/configurations/CardFeatureFlag",
            resolve
        );
    });

    let configData = await new Promise((resolve, reject) => {
        let data = '';
        res.on('data', chunk => data += chunk);
        res.on('error', err => reject(err));
        res.on('end', () => resolve(data));
    });

    const parsedConfigData = JSON.parse(configData);

    const DynamoParams = {
        TableName: 'AWSCribsRentalMansions'
    };
    
//Fetching the listings from DynamoDB

    async function listItems() {
        try {
            const data = await docClient.scan(DynamoParams).promise();
            return data;
        } catch (err) {
            return err;
        }
    }

//Checking for the Caraousel Feature Flag
    if (parsedConfigData.showcarousel.enabled == true) {
        let returnhtml = ``;

        try {

            const data = await listItems();

            for (let i = 0; i < parsedConfigData.pagination.number; i++) {
                returnhtml += `<div class="col-md-4 mt-4">
    		    <div class="card profile-card-5">
    		        <div class="card-img-block">

                        <div class="slideshow-container">`;

                for (let j = 0; j < data.Items[i].Image.length; j++) {
                    returnhtml += `<div class="mySlides` + (i + 1) + `">
                    <img class="card-img-top" style="height: 300px;" src="` + data.Items[i].Image[j].name + `" style="width:100%">
                  </div>`;
                }

                returnhtml += `</div>`;
                
                if(data.Items[i].Image.length > 1) {
                    returnhtml += `<a class="prev" onclick="plusSlides(-1, ` + i + `)">&#10094;</a>
                            <a class="next" onclick="plusSlides(1, ` + i + `)">&#10095;</a>`;
                }
                
    		        returnhtml += `</div>
                    <div class="card-body pt-0">
                    <h5 class="card-title">` + data.Items[i].Name + ` <span style="font-size: 0.7em;color:rgb(255, 64, 64)">(` + data.Items[i].Location + `)</span></h5>
                    <p class="card-text">` + data.Items[i].Description + `</p>
                    <a class="btn btn-primary"style="display: inline" href="#">Check Availability</a>
                    <span style="float: right;cursor: pointer;" onclick="favoriteStar(this)"><span class="fa fa-star"></span></span>
                  </div>
                </div>
    		</div>`;
            }

            return {
                statusCode: 200,
                body: returnhtml,
            };

        } catch (err) {
            return {
                error: err
            }
        }

    } else {

        let returnhtml = ``;

        try {

            const data = await listItems();
            console.log("dynamo db data: ", data)

//Checking for Pagination Numbers
            for (let i = 0; i < parsedConfigData.pagination.number; i++) {
                returnhtml += `<div class="col-md-4 mt-4">
    		    <div class="card profile-card-5">
    		        <div class="card-img-block">
                    <img class="card-img-top" style="height: 300px;" src="` + data.Items[i].Image[0].name + `" style="width:100%"
                        alt="Card image cap" style="height: 300px;">
    		        </div>
                    <div class="card-body pt-0">
                    <h5 class="card-title">` + data.Items[i].Name + ` <span style="font-size: 0.7em;color:rgb(255, 64, 64)">(` + data.Items[i].Location + `)</span></h5>
                    <p class="card-text">` + data.Items[i].Description + `</p>
                    <a class="btn btn-primary"style="display: inline" href="#">Check Availability</a>
                    <span style="float: right;cursor: pointer;" onclick="favoriteStar(this)"><span class="fa fa-star"></span></span>
                  </div>
                </div>
    		</div>`;
            }

            return {
                statusCode: 200,
                body: returnhtml,
            };

        } catch (err) {
            return {
                error: err
            };
        }
    }
};

