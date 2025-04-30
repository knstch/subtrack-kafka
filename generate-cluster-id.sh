set -euo pipefail
ENV_FILE=".env"

KAFKA_IMAGE="confluentinc/cp-kafka:7.5.0"

if [[ -f "$ENV_FILE" ]] && grep -q "^CLUSTER_ID=" "$ENV_FILE"; then
  echo "CLUSTER_ID already exists $ENV_FILE – skip."
  exit 0
fi

echo "Generating new CLUSTER_ID…"
CLUSTER_ID=$(docker run --rm "$KAFKA_IMAGE" kafka-storage random-uuid)

{
  # если .env не пустой и не заканчивается \n – добавим пустую строку
  [ -s "$ENV_FILE" ] && tail -c1 "$ENV_FILE" | read -r _ || echo
  echo "CLUSTER_ID=$CLUSTER_ID"
} >> "$ENV_FILE"
echo "Written in $ENV_FILE"