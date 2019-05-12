let fanSpeed = $("#fan_speed");
fanSpeed.slider({value:speed,});

$(".slider-horizontal").mouseup(function () {
    $.ajax({
        url: '/changeSpeed',
        type: "post",
        dataType: 'html',
        data: { speed: fanSpeed.val()} ,
        success: function (data) {
            console.log(data);
        }
    });
})