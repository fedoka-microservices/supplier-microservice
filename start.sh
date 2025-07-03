#!/bin/sh

echo "Esperando a que la base de datos esté lista..."

# Intentamos conectarnos cada 2 segundos, hasta un máximo de 30 intentos
for i in $(seq 1 30); do
  nc -z "$DB_HOST" "$DB_PORT" && break
  echo "Esperando..."
  sleep 2
done

# Si después de 30 intentos sigue sin funcionar, salimos
if ! nc -z "$DB_HOST" "$DB_PORT"; then
  echo "No se pudo conectar a la base de datos en $DB_HOST:$DB_PORT"
  exit 1
fi

echo "Base de datos lista. Ejecutando migraciones..."

# Ejecutar migraciones con goose
goose -dir db/migrations mysql "$DB_USERNAME:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME" up

echo "Iniciando servicio..."
exec "$@"
