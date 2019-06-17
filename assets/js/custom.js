

$( document ).ready(function() {

    var fanSpeed = $("#fan_speed");
    var originalVal;

    fanSpeed.slider({value:speed,});
    fanSpeed.slider().on('slideStart', function(ev){
    	originalVal = fanSpeed.data('slider').getValue();
    });

    fanSpeed.slider().on('slideStop', function(ev){
    	var newVal = fanSpeed.data('slider').getValue();
    	if(originalVal != newVal) {    		
    		$.ajax({
    			url: '/changeSpeed',
    			type: "post",
    			dataType: 'html',
    			data: { speed: newVal} ,
    			success: function (data) {
    				console.log(data);
    			}
    		});
    	}
    });


    var ajaxInterval = 60000;
    setInterval(function(){ 
        $.ajax({
            url: '/updateData',
            type: "get",
            success: function (data) {
                var upperTemperature = $('.temperature_Upper');
                upperTemperature.html(data.SensorsData.Temperature.current["0"].Value + ' °C');
                var lowerTemperature = $('.temperature_Lower');
                lowerTemperature.html(data.SensorsData.Temperature.current["1"].Value + ' °C');

                var upperHumidity = $('.humidity_Upper');
                upperHumidity.html(data.SensorsData.Humidity.current["0"].Value + ' %'); 
                var lowerHumidity = $('.humidity_Lower');
                lowerHumidity.html(data.SensorsData.Humidity.current["1"].Value + ' %');

                var soilSurface = $('.soil_Surface');
                soilSurface.html(data.SensorsData.SoilMoisture.current["0"].Value + ' %');
                var soilRoot = $('.soil_Root');
                soilRoot.html(data.SensorsData.SoilMoisture.current["1"].Value + ' %');
                
                console.log(data);
            }
        });
    }, ajaxInterval);
    
});

let ctx = document.getElementById('myChart').getContext('2d');
let chart = new Chart(ctx, {
    // The type of chart we want to create
    type: 'line',

    // The data for our dataset
    data: {
        labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
        datasets: [
            {
            label: 'My First dataset',
            // backgroundColor: 'rgb(255, 99, 132)',
            borderColor: 'rgb(255, 99, 132)',
            data: [0, 10, 5, 2, 20, 30, 45]
            },
            {
                label: 'My second dataset',
                // backgroundColor: 'rgb(255, 99, 132)',
                borderColor: 'rgb(255, 99, 132)',
                data: [10, 20, 5, 2, 10, 20, 25]
            }
        ]
    },

    // Configuration options go here
    options: {}
});