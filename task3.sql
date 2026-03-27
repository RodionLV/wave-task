-- PostgreSQL --

CREATE TABLE IF NOT EXISTS devices (
    id			SERIAL PRIMARY KEY,
    hostname	VARCHAR(70) NOT NULL,
    ip 			INET NOT NULL,
    location	TEXT NOT NULL,
    is_active	BOOL DEFAULT true,
    created_at	TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE IF NOT EXISTS configs (
    id SERIAL PRIMARY KEY,

    -- other params
    some_data TEXT,

    device_id INT,
    CONSTRAINT fk_device FOREIGN KEY(device_id) 
        REFERENCES devices(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS logs (
    id SERIAL PRIMARY KEY,

    message TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),

    device_id INT,
    CONSTRAINT fk_device FOREIGN KEY(device_id) 
        REFERENCES devices(id) ON DELETE CASCADE
);

INSERT INTO devices(hostname, ip, location) VALUES 
('host1', '123.123.123.121', 'city1'),
('host2', '123.123.123.122', 'city2'),
('host3', '123.123.123.123', 'city3'),
('host4', '123.123.123.124', 'city4'),
('host5', '123.123.123.125', 'city5'),
('host6', '123.123.123.126', 'city6'),
('host7', '123.123.123.127', 'city7'),
('host8', '123.123.123.128', 'city8');

INSERT INTO configs(device_id, some_data) VALUES
(1, 'data1'),
(1, 'data2'),
(2, 'data3'),
(2, 'data4'),
(2, 'data5'),
(3, 'data6'),
(4, 'data7'),
(5, 'data8'),
(6, 'data9');

INSERT INTO logs(device_id, message) VALUES
(1, 'data1'),
(1, 'data2'),
(2, 'data3'),
(1, 'data4'),
(2, 'data5'),
(3, 'data6'),
(1, 'data7'),
(5, 'data8'),
(6, 'data9');



-- test sql query 1
SELECT d.id, d.hostname, d.ip, d.location, COUNT(c.device_id) as count_configs 
    FROM devices d 
    LEFT JOIN configs c ON d.id = c.device_id 
    WHERE d.is_active = true
    GROUP BY d.id;


CREATE INDEX configs_index ON configs(device_id);


-- test sql query 2

-- $id = 1, $N = 10

SELECT * FROM logs l WHERE l.device_id = 1 ORDER BY created_at DESC LIMIT 10;
 

CREATE INDEX logs_create_at_index ON logs(created_at);
CREATE INDEX logs_device_id_index ON logs USING HASH(device_id);


