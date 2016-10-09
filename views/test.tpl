<!doctype html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport"
              content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
        <meta http-equiv="X-UA-Compatible" content="ie=edge">
        <title>Test map</title>
        <link rel="stylesheet" href="/static/css/leaflet.css">
        <style>
            body, html {
                height: 100%;
                padding: 0;
                margin: 0;
            }

            .control{
                z-index: 999;
                position: fixed;
                //width: 50%;
                top: 0;
                right:0;
                background: #f4f4f4;
                padding: 10px;
                border-left: 1px solid #cccccc;
                border-bottom: 1px solid #cccccc;
            }

            #map {
                height: 600px;
            }
            #select-map{
                margin-right: 20px;
                min-width: 200px;
            }
        </style>
    </head>
    <body>

        <div class="control" style="he">
            <select name="select-map" id="select-map">
                <option value="none">none</option>
            </select>

            <label for="check-gs">
                Use gray-scale tile convert
                <input type="checkbox" name="check-gs" id="check-gs">
            </label>
        </div>

        <div id="map" style="height: 100%;"></div>

        <script src="/static/js/jquery-3.1.1.min.js"></script>
        <script src="/static/js/leaflet.js"></script>
        <script>
            (function () {
                function updateTileLayer() {

                    //Remove all layers
                    map.eachLayer(function (layer) {
                        map.removeLayer(layer);
                    });

                    //Add new layer
                    L.tileLayer('http://localhost:8080/?x={x}&y={y}&z={z}&server=' + serverAlias + (useGS ? "&gs=1" : ""), {
                        maxZoom: 18,
                    }).addTo(map);
                }

                var serversConfig = JSON.parse({{.serversConfig}});

                var map = L.map('map').setView([56.132326, 40.404282], 13);

                var useGS = false;
                var serverAlias = "";

                $.each(serversConfig['servers'], function (key, server) {
                    $('#select-map')
                            .append($('<option></option>')
                                    .attr('value', server.alias)
                                    .text(server.alias));
                });

                $('#select-map').on('change', function () {
                    serverAlias = $(this).val();
                    updateTileLayer();
                });

                $('#check-gs').on('change', function () {
                    useGS = $(this).prop("checked");
                    updateTileLayer();
                });
            })()
        </script>
    </body>
</html>