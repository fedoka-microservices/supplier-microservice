-- +goose Up
INSERT INTO suppliers (name, email, phone, address, rfc, created_at, updated_at) VALUES
('VALMEX', 'ventas@valmex.mx', '9612345672', 'Calle Sur 45', 'VAL020304BC2', NOW(), NOW()),
('AUTOREPUESTOS JIMÉNEZ', 'contacto@jimenez.mx', '9612233445', 'Avenida Central 100', 'JIM010203DF1', NOW(), NOW());

-- +goose Down
DELETE FROM suppliers WHERE rfc IN ('VAL020304BC2', 'JIM010203DF1');
