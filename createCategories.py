from pymilvus import (
    connections,
    FieldSchema, CollectionSchema, DataType,
    Collection, utility
)
import time

# Connect to Milvus
connections.connect(
    alias="default",
    host="localhost",  # Replace with your Milvus host if different
    port="19530"
)

def create_resumes_collection():
    collection_name = "resumes"

    # Check if the collection already exists
    if utility.has_collection(collection_name):
        print(f"Collection '{collection_name}' already exists. Dropping it for recreation.")
        utility.drop_collection(collection_name)
        print(f"Collection '{collection_name}' dropped.")

    # Define the fields in the collection
    fields = [
        FieldSchema(name="resume_id", dtype=DataType.INT64, is_primary=True, auto_id=False),
        FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=1024)
    ]

    # Create a schema for the collection
    schema = CollectionSchema(fields, description="Resume embeddings")

    # Create the collection
    collection = Collection(name=collection_name, schema=schema)
    print(f"Collection '{collection_name}' created successfully.")

    # Create an index on the 'embedding' field
    index_params = {
        "metric_type": "L2",
        "index_type": "IVF_FLAT",
        "params": {"nlist": 128},
    }

    print(f"Creating index on 'embedding' field...")
    collection.create_index(field_name="embedding", index_params=index_params)
    print(f"Index created successfully.")

    # Load the collection into memory
    print(f"Loading collection '{collection_name}' into memory...")
    collection.load()
    print(f"Collection '{collection_name}' loaded into memory.")

    # Optionally, insert some dummy data to ensure everything works
    # insert_dummy_data(collection)

def insert_dummy_data(collection):
    # Generate some dummy data
    resume_ids = [1, 2, 3]
    embeddings = [[float(i) for i in range(1024)] for _ in range(3)]

    # Prepare data in a dictionary
    data = {
        "resume_id": resume_ids,
        "embedding": embeddings
    }

    # Insert data into the collection
    print(f"Inserting dummy data into collection '{collection.name}'...")
    mr = collection.insert(data)
    print(f"Inserted {len(mr.primary_keys)} records into collection '{collection.name}'.")

    # Flush to ensure data is persisted
    print(f"Flushing collection '{collection.name}'...")
    utility.flush([collection.name])
    print(f"Flush completed.")

if __name__ == "__main__":
    create_resumes_collection()
    # Uncomment the next line if you want to insert dummy data
    # insert_dummy_data(Collection("resumes"))