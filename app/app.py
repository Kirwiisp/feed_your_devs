import os
import discord
import logging
from responses import get_response
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

async def send_message(message: discord.Message, user_mess: str) -> None:
    if not user_mess:
        print("Message empty. Might be intents misconfiguration")
        return

    is_private = user_mess[0] == '?'
    if is_private: 
        user_mess = user_mess[1:]
    
    try:
        response: str = get_response(user_mess)
        await message.author.send(response) if is_private else await message.channel.send(response)
    except Exception as e:
        logging.error("Error in sending message")

@client.event
async def on_ready():
    #for guild in client.guilds:
    #    print(f"Guild = ",guild.name)
    #    members = guild.members
    #    for member in members:
    #        print(f"    User: ",member.name)
    print(f'{client.user.name} has connected to Discord!')

@client.event
async def on_message(message: discord.Message) -> None:
    if message.author == client.user:
        return
    username: str = str(message.author)
    user_message: str = str(message.content)
    channel: str = str(message.channel)

    print(f'[{channel}] {username}: "{user_message}"')
    await send_message(message, user_message)

def main() -> None:
# Run the bot
    client.run(TOKEN)
if __name__ == '__main__':
    main()

