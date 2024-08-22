from confluent_kafka import Consumer, KafkaError
from collections import defaultdict

# Kafka configuration
conf = {
    'bootstrap.servers': "localhost:9092",
    'group.id': "search-query-analytics",
    'auto.offset.reset': 'earliest'
}

# Create Consumer instance
consumer = Consumer(**conf)

# Subscribe to the topic
topic = "search-queries"
consumer.subscribe([topic])

# Dictionary to store search query frequencies
query_counts = defaultdict(int)

print("Listening for search queries...")

try:
    while True:
        msg = consumer.poll(timeout=1.0)  # Wait for a message

        if msg is None:
            continue

        if msg.error():
            if msg.error().code() == KafkaError._PARTITION_EOF:
                continue
            else:
                print(msg.error())
                break

        # Get the search query from the message
        search_query = msg.value().decode('utf-8')

        # Update the count for the search query
        query_counts[search_query] += 1

        # Print the current top search queries
        print("\nTop Search Queries:")
        for query, count in sorted(query_counts.items(), key=lambda item: item[1], reverse=True):
            print(f"{query}: {count} times")

except KeyboardInterrupt:
    pass
finally:
    # Close down consumer to commit final offsets.
    consumer.close()
