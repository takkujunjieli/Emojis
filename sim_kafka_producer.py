from confluent_kafka import Producer
import time
import random

# Kafka configuration
conf = {'bootstrap.servers': "localhost:9092"}

# Create Producer instance
producer = Producer(**conf)

# List of sample search queries
search_queries = ["Kafka tutorial", "Python Kafka", "Real-time analytics", 
                  "Kafka Streams", "Distributed systems", "Big data", "Python streams"]

def delivery_report(err, msg):
    if err is not None:
        print('Message delivery failed: {}'.format(err))
    else:
        print('Message delivered to {} [{}]'.format(msg.topic(), msg.partition()))

# Send search queries to the Kafka topic
topic = "search-queries"
for _ in range(100):  # Simulating 100 search queries
    query = random.choice(search_queries)
    producer.produce(topic, key="query", value=query, callback=delivery_report)
    producer.poll(0)  # Serve delivery reports
    time.sleep(1)  # Simulate a delay between search queries

# Wait for any outstanding messages to be delivered and delivery reports received
producer.flush()
