import os
import discord
from dotenv import load_dotenv

# Load environment variables
load_dotenv()
TOKEN = os.getenv("TOKEN")

# Set up intents
intents = discord.Intents.default()  # Use default intents
intents.messages = True  # Enable intents for messages 
intents.members = True # Enable Intents for members

# Create the client with intents
client = discord.Client(intents=intents)
@client.event
async def on_ready():
    for guild in client.guilds:
        print(f"Guild = ",guild.name)
        members = guild.members
        for member in members:
            print(f"User: ",member.name)
    print(f'{client.user.name} has connected to Discord!')

# Run the bot
client.run(TOKEN)

