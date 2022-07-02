fetch('https://xe5phfcmjxvg7fdy4nqvrop62m0nvdgm.lambda-url.us-east-1.on.aws/')
    .then(response => response.text())
    .then((data) => {
        document.getElementById("displayCecil").innerHTML = data;
        console.log(data);
    });