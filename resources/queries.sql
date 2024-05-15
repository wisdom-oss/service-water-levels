-- name: get-measurement-stations
SELECT *
FROM geodata.water_level_stations;

-- name: get-single-stations
SELECT *
FROM geodata.water_level_stations
WHERE website_id = $1;