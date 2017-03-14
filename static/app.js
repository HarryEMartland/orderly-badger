$(function () {

    var containerTableBody = $('#containerTableBody');

    function loadContainers() {
        $.get("/containers", function (data) {
            var now = moment();
            containerTableBody.empty();
            $.each(JSON.parse(data), function (i, container) {
                appendContainer(container, now)
            })
        });
    }

    function appendContainer(container, now) {
        var timeMiliseconds = container.startedAt * 1000;

        containerTableBody.append("<tr class='container-row' data-container-id='" + container.id + "'>" +
            "<td title='" + container.id + "'>" + container.name + "</td>" +
            "<td>" + moment(timeMiliseconds).format("DD/MM/YY HH:MM") + "</td>" +
            "<td>" + moment.duration(container.maxAge, 'seconds').humanize() + "</td>" +
            "<td class='timeLeft' data-started-at='" + timeMiliseconds + "' data-max-age='" + container.maxAge + "'>" +
            calculateTimeLeft(timeMiliseconds, container.maxAge, now)
            + "</td>" +
            "</tr>")
    }

    setInterval(function () {
        var now = moment();
        $(".timeLeft").each(function (i, element) {
            element = $(element);
            var startedAt = parseInt(element.attr("data-started-at"));
            var maxAge = parseInt(element.attr("data-max-age"));
            element.text(calculateTimeLeft(startedAt, maxAge, now))
        })
    }, 1000);

    function calculateTimeLeft(startedAt, maxAge, now) {
        var startedAt = moment((startedAt));
        var maxAge = moment.duration((maxAge), 'seconds');

        var deleteTime = startedAt.add(maxAge);

        return moment.duration(deleteTime.diff(now)).format("h [hrs], m [min], s [sec]");
    }


    var ws = new ReconnectingWebSocket("ws://" + window.location.host + "/ws");

    ws.onopen = function () {
        loadContainers()
    };

    ws.onmessage = function (evt) {
        var message = JSON.parse(evt.data);

        if (message.type === "start") {
            appendContainer(message.data, moment())
        }
        if (message.type === "die") {
            $(".container-row[data-container-id='" + message.data + "']").remove()
        }
    };

    ws.onclose = function () {
        console.error("Connection is closed...");
    };
});