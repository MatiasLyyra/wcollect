CREATE DATABASE weather;

CREATE DICTIONARY weather.smart_symbol_dict
(
    Id Int32,
    Description String
)
PRIMARY KEY Id
SOURCE(FILE(
    path '/var/lib/clickhouse/user_files/smart_symbols.csv'
    format 'CSVWithNames'
))
LIFETIME(3600)
LAYOUT(FLAT());

CREATE TABLE weather.observations
(
    Location                LowCardinality(String),
    ObservationTime         DateTime('UTC'),
    CreationTime            DateTime('UTC'),
    Pressure                Float32 CODEC(Gorilla),
    PrecipitationAmount     Float32 CODEC(Gorilla),
    PrecipitationIntensity  Float32 CODEC(Gorilla),
    RelativeHumidity        Float32 CODEC(Gorilla),
    SnowDepth               Float32 CODEC(Gorilla),
    AirTemperature          Float32 CODEC(Gorilla),
    Dewpoint                Float32 CODEC(Gorilla),
    Visibility              Float32 CODEC(Gorilla),
    WindDirection           Float32 CODEC(Gorilla),
    GustSpeed               Float32 CODEC(Gorilla),
    WindSpeed               Float32 CODEC(Gorilla),
    SmartSymbol             Int32 CODEC(ZSTD),
    WeatherDescription      String MATERIALIZED dictGet('weather.smart_symbol_dict', 'Description', IF(SmartSymbol >= 100, SmartSymbol - 100, SmartSymbol))
)
ENGINE = ReplacingMergeTree(CreationTime)
ORDER BY (Location, ObservationTime)
PARTITION BY (toStartOfMonth(ObservationTime));

CREATE TABLE weather.forecasts
(
    Location                LowCardinality(String),
    ObservationTime         DateTime('UTC'),
    CreationTime            DateTime('UTC'),
    Pressure                Float32 CODEC(Gorilla),
    PrecipitationAmount     Float32 CODEC(Gorilla),
    RelativeHumidity        Float32 CODEC(Gorilla),
    AirTemperature          Float32 CODEC(Gorilla),
    Dewpoint                Float32 CODEC(Gorilla),
    WindDirection           Float32 CODEC(Gorilla),
    GustSpeed               Float32 CODEC(Gorilla),
    WindSpeed               Float32 CODEC(Gorilla),
    SmartSymbol             Int32 CODEC(ZSTD),
    WeatherDescription      String MATERIALIZED dictGet('weather.smart_symbol_dict', 'Description', IF(SmartSymbol >= 100, SmartSymbol - 100, SmartSymbol))
)
ENGINE = ReplacingMergeTree(CreationTime)
ORDER BY (Location, ObservationTime)
PARTITION BY (toStartOfMonth(ObservationTime));

CREATE VIEW weather.latest_observations AS
SELECT
    Location,
    ObservationTime,
    argMax(Pressure, CreationTime) AS Pressure,
    argMax(PrecipitationAmount, CreationTime) AS PrecipitationAmount,
    argMax(PrecipitationIntensity, CreationTime) AS PrecipitationIntensity,
    argMax(RelativeHumidity, CreationTime) AS RelativeHumidity,
    argMax(SnowDepth, CreationTime) AS SnowDepth,
    argMax(AirTemperature, CreationTime) AS AirTemperature,
    argMax(Dewpoint, CreationTime) AS Dewpoint,
    argMax(Visibility, CreationTime) AS Visibility,
    argMax(WindDirection, CreationTime) AS WindDirection,
    argMax(GustSpeed, CreationTime) AS GustSpeed,
    argMax(WindSpeed, CreationTime) AS WindSpeed,
    argMax(SmartSymbol, CreationTime) AS SmartSymbol,
    argMax(WeatherDescription, CreationTime) AS WeatherDescription
FROM weather.observations
GROUP BY Location, ObservationTime;

CREATE VIEW weather.weather_combined AS
SELECT
    Location,
    toStartOfInterval(ObservationTime, INTERVAL 1 HOUR) as ObservationTime,
    avg(Pressure) AS Pressure,
    avg(PrecipitationAmount) AS PrecipitationAmount,
    avg(RelativeHumidity) AS RelativeHumidity,
    avg(AirTemperature) AS AirTemperature,
    avg(Dewpoint) AS Dewpoint,
    avg(WindDirection) AS WindDirection,
    avg(GustSpeed) AS GustSpeed,
    avg(WindSpeed) AS WindSpeed
FROM weather.latest_observations
GROUP BY Location, ObservationTime
UNION ALL
SELECT
    Location,
    toStartOfInterval(ObservationTime, INTERVAL 1 HOUR) as ObservationTime,
    avg(Pressure) AS Pressure,
    avg(PrecipitationAmount) AS PrecipitationAmount,
    avg(RelativeHumidity) AS RelativeHumidity,
    avg(AirTemperature) AS AirTemperature,
    avg(Dewpoint) AS Dewpoint,
    avg(WindDirection) AS WindDirection,
    avg(GustSpeed) AS GustSpeed,
    avg(WindSpeed) AS WindSpeed
FROM weather.latest_forecasts
GROUP BY Location, ObservationTime;



CREATE VIEW weather.latest_forecasts AS
SELECT
    Location,
    ObservationTime,
    argMax(Pressure, CreationTime) AS Pressure,
    argMax(PrecipitationAmount, CreationTime) AS PrecipitationAmount,
    argMax(RelativeHumidity, CreationTime) AS RelativeHumidity,
    argMax(AirTemperature, CreationTime) AS AirTemperature,
    argMax(Dewpoint, CreationTime) AS Dewpoint,
    argMax(WindDirection, CreationTime) AS WindDirection,
    argMax(GustSpeed, CreationTime) AS GustSpeed,
    argMax(WindSpeed, CreationTime) AS WindSpeed,
    argMax(SmartSymbol, CreationTime) AS SmartSymbol,
    argMax(WeatherDescription, CreationTime) AS WeatherDescription
FROM weather.forecasts
GROUP BY Location, ObservationTime;


CREATE USER weather_writer IDENTIFIED WITH bcrypt_password BY '<password>';
CREATE USER weather_reader IDENTIFIED WITH bcrypt_password BY '<password>';
GRANT INSERT ON weather.observations TO weather_writer;
GRANT INSERT ON weather.forecasts TO weather_writer;
GRANT SELECT ON weather.* TO weather_reader;
GRANT dictGet ON weather.smart_symbol_dict TO weather_writer;