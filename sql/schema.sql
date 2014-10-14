CREATE TABLE sensor_controller (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    mac_address CHAR(12) NOT NULL,
    last_startup DATETIME NOT NULL,
    CONSTRAINT pk_sensor_controller_id PRIMARY KEY (id),
    CONSTRAINT un_mac_address UNIQUE INDEX (mac_address)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8;

CREATE TABLE sensor (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    device_address CHAR(16) NOT NULL,
    CONSTRAINT pk_sensor_id PRIMARY KEY (id),
    CONSTRAINT un_device_address UNIQUE INDEX (device_address)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8;

CREATE TABLE sensor_temperature (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    sensor_id INT UNSIGNED NOT NULL,
    timestamp DATETIME NOT NULL,
    value DECIMAL(5, 2) NOT NULL,
    CONSTRAINT pk_sensor_temperature_id PRIMARY KEY (id),
    INDEX id_sensor_id (sensor_id),
    CONSTRAINT fk_sensor_id FOREIGN KEY (sensor_id) REFERENCES sensor (id)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8;
