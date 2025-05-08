from pyspark.sql import SparkSession
from urllib.parse import urlparse

import boto3
import sys

# get bucket and key
parsed_uri = urlparse(sys.argv[1])

# Extract bucket name and key
bucket_name = parsed_uri.netloc
file_key = parsed_uri.path.lstrip('/')

# download file
s3 = boto3.client('s3')
response = s3.get_object(Bucket=bucket_name, Key=file_key)
query = response['Body'].read().decode('utf-8')

# execute query, write result in avro format to S3
spark = SparkSession.builder.appName("spark-sql-app").getOrCreate()
spark.sql(query).write.format("avro").save(sys.argv[2])
