import os
import discord
from dotenv import load_dotenv

# Load environment variables
load_dotenv()
TOKEN = os.getenv("TOKEN")

# Set up intents
intents = discord.Intents.default()  # Use default intents
intents.messages = True  # Enable intents for messages (adjust based on your needs)

# Create the client with intents
client = discord.Client(intents=intents)

@client.event
async def on_ready():
    print(f'{client.user} has connected to Discord!')

# Run the bot
client.run(TOKEN)

