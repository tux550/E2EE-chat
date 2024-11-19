# Use official Python image as a base
FROM python:3.11-slim

# Set the working directory in the container
WORKDIR /app

# Install dependencies
RUN pip install --upgrade pip
COPY requirements.txt .
RUN pip install -r requirements.txt

# Copy the application files into the container
COPY . .

# Expose the port that FastAPI will run on
EXPOSE 8081 8082

# Run the application with uvicorn
CMD ["python", "api-gateway.py"]