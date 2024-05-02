import psycopg2
import numpy as np
import matplotlib.pyplot as plt

connection = psycopg2.connect(database="healthtracker", user="healthtracker", password="12345", host="localhost", port=5432)

cursor = connection.cursor()

cursor.execute("SELECT date(created_at), calories from healthtracker where user_id = 2")

record = cursor.fetchall()

time = []
calories = []

for i in record:
    time.append(i[0])
    calories.append(i[1])

plt.bar(time, calories)
plt.ylim(0, 3000)
plt.xlabel("Time")
plt.ylabel("Calories")
plt.title("Information")

plt.show()