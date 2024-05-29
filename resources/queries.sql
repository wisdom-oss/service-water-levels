-- name: get-measurement-stations
SELECT *
FROM geodata.water_level_stations;

-- name: get-single-station
SELECT *
FROM geodata.water_level_stations
WHERE website_id = $1;

-- name: get-measurements
SELECT *
FROM timeseries.nlwkn_water_levels
ORDER BY station, date;

-- name: get-measurements-for-station
SELECT *
FROM timeseries.nlwkn_water_levels
WHERE station = $1
ORDER BY station, date;

-- name: get-measurements-in-range
SELECT *
FROM timeseries.nlwkn_water_levels
WHERE date BETWEEN $1 AND $2
ORDER BY station, date;

-- name: get-measurements-for-station-in-range
SELECT *
FROM timeseries.nlwkn_water_levels
WHERE date BETWEEN $1::date AND $2::date
AND station = $3
ORDER BY date;