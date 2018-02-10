$(function () {
    'use strict';

    function makeRequest(opts) {
        return new Promise(function (resolve, reject) {
            var xhr = new XMLHttpRequest();
            xhr.open(opts.method, opts.url);
            xhr.onload = function () {
                if (this.status >= 200 && this.status < 300) {
                    var payload = xhr.responseText;
                    if (xhr.getResponseHeader("Content-Type")==="application/json") {
                        payload = JSON.parse(xhr.responseText);
                        resolve(payload);
                    }
                    // resolve(xhr.response);
                    resolve(payload);
                } else {
                    reject({
                        status: this.status,
                        statusText: xhr.statusText
                    });
                }
            };
            xhr.onerror = function () {
                reject({
                    status: this.status,
                    statusText: xhr.statusText
                });
            };
            if (opts.headers) {
                Object.keys(opts.headers).forEach(function (key) {
                    xhr.setRequestHeader(key, opts.headers[key]);
                });
            }
            var params = opts.params;
            var payload;
            // We'll need to stringify if we've been given an object
            // If we have a string, this is skipped.
            if (params && typeof params === 'object') {
                payload = Object.keys(params).map(function (key) {
                    return encodeURIComponent(key) + '=' + encodeURIComponent(params[key]);
                }).join('&');
            }
            var json = opts.json;
            if (json && typeof json === 'object') {
                xhr.setRequestHeader("Content-type", "application/json");
                payload = JSON.stringify(json);
            }
            xhr.send(payload);
        });
    }

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

    var template;
    makeRequest({method: 'GET', url: '/static/tmpl/post.html'}).then(function(tmpl){
        template = Handlebars.compile(tmpl);

        if (!window['WebSocket']) {
            console.error('Błąd: Przeglądarka nie wspiera WebSocket')
        } else {
            connect();
        }
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
            var d = template(data);
            console.log(data);
            messages.append(d);
        };
    }
});
