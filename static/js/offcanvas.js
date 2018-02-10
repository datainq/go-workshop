$(function () {
    'use strict'

    $('[data-toggle="offcanvas"]').on('click', function () {
        $('.offcanvas-collapse').toggleClass('open')
    });

    var socket = null;
    var msgBox = $('#user-msg-box-ta');
    var messages = $('#messages');
    $('#user-msg-box').submit(function() {
        if (!msgBox.val()) return false;
        if (!socket) {
            alert('Błąd: brak połączenia z serwerem.');
            return false;
        }
        var data = {author: 'Pawel', Message: msgBox.val()};
        socket.send(JSON.stringify(data));
        msgBox.val('');
        return false;
    });

    function connect() {
        socket = new WebSocket("ws://localhost:8080/ws");
        var status = socket.CONNECTING;
        socket.onopen = function (ev) {
            console.log('websocket opened');
            status = WebSocket.OPEN;
        };
        socket.onerror = function (ev) {
            console.log('problem opening a websocket');
        };
        socket.onclose = function (ev) {
            if (status===WebSocket.OPEN) {
                console.error("Połączenie zerwane.");
            }
            setTimeout(connect, 5);
            status = WebSocket.CLOSED;
        };
        socket.onmessage = function (ev) {
            var data = JSON.parse(ev.data);
            var d = document.createTextNode(data.Message);
            console.log(data);
            messages.append(d);
        };
    }

    if (!window['WebSocket']) {
        console.error('Błąd: Przeglądarka nie wspiera WebSocket')
    } else {
        connect();
    }
});
