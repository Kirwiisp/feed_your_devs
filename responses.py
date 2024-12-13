from random import choice, randint
import re   

def get_response(user_input: str) -> str:
    lowered: str= user_input.lower()

    if lowered == '':
        return "Message missing"
    elif 'help' in lowered:
        return 'List of commands : [...]'
    elif '/r' in lowered:
        match = re.findall(r"(?<=^/r)\s*([0-9]+)d([0-9]+)", lowered)
        if match:
            numb = int(match[0][0])
            die = int(match [0][1])
            total = sum(randint(1,die) for _ in range(numb))
            return f"Result : {total}"
        else:
            return "Invalid dice roll format. Use /r <numberOfDice>d<DiceValue> (e.g. '/r 1d6'"
    else:
        return choice(["Couldn't understand your message, please try again","I didn't understand","What ? Can you repeat?"])
