import requests
import time
import yaml, json



'''
check if a card is on the battlefield
:param cardPos: TopLeftY of card, INT
:param screenHeight: height of play screen INT
:param screenWidth: width of play screen INT
:return: True if card is at a played Y position, False if not at a played Y position
'''
def inPlay(cardPos, screenHeight, screenWidth):
    inPlayLocal = round(screenHeight * 0.41666666666)
    inPlayOpp = round(screenHeight * 0.73148148148)

    if cardPos == inPlayLocal or cardPos == inPlayOpp:
        return True
    else:
        return False
'''
check which side is attacking by which side plays cards onto the battlefield first
:param allCards: list of rectangle dictionaries with cardID, CardCode,top left X and Y, Card width and Height, and which player has it
:param height: height of play screen INT
:param width: width of play screen INT
:return: 1 if attacker should be local player, -1 if attacker is the opponent, 0 if unsure (both sides have no play, both sides populated, etc)
'''
def attackerCheck(allCards, height, width):
    local = []
    opponent = []
    for rectangle in allCards:
        #local player card in play
        if rectangle['LocalPlayer'] == True and inPlay(rectangle['TopLeftY'], height, width):
            local.append(rectangle['CardID'])
        #opponent card in play
        elif rectangle['LocalPlayer'] == False and inPlay(rectangle['TopLeftY'], height, width):
            opponent.append(rectangle['CardID'])
    if len(local) > 0 and len(opponent) == 0:
        return 1
    elif len(opponent) > 0 and len(local) == 0:
        return -1
    else:
        return 0
'''
creates set of cards being played at that moment
:param tempSet: the current tempSet of cards
:param playedSet: the current cards in play on this poll
:return: new tempSet of cards being played right now
'''
def playedCardsChooser(tempSet, playedSet):
    newSet = set()
    for card in playedSet:
        if card not in tempSet:
            newSet.add(card)
        if card in tempSet:
            newSet.add(card)
    return newSet


def main():
    gamesPlayed = -1
    gameState = 0
    pastGameState = 0
    tempLocal = set()
    cardsPlayedLocal = set()
    id = ''

    print("Program Started")
    while True:
        #only poll once a second
        time.sleep(1)
        #get card position and game result data
        try:
            req = requests.get("http://localhost:21337/positional-rectangles")
            reqGame = requests.get('http://localhost:21337/game-result')
        except:
            print('Lost Connection To Legends of Runeterra \n Closing War on Runeterra Client')
            break

        if req.status_code == 200:
            received = yaml.safe_load(req.text)
            #print(received)
            if received['GameState'] == 'InProgress':
                #find and set attacker states
                attack = attackerCheck(received['Rectangles'], received['Screen']['ScreenHeight'], received['Screen']['ScreenWidth'])
                if attack == 1:
                    pastGameState = gameState
                    gameState = attack
                elif attack == -1:
                    pastGameState = gameState
                    gameState = attack

                #populate tempSets with cards that potentially are played
                played = set()
                for rectangle in received['Rectangles']:
                    #if attacker is local player
                    if gameState == 1:
                        #potential attack, can be changed
                        if rectangle['LocalPlayer'] == True and inPlay(rectangle['TopLeftY'], received['Screen']['ScreenHeight'], received['Screen']['ScreenWidth']):
                            played.add(rectangle['CardCode'])
                        #if opponent defense then attack is locked and these cards are confirmed played
                        if rectangle['LocalPlayer'] == True and inPlay(rectangle['TopLeftY'], received['Screen']['ScreenHeight'], received['Screen']['ScreenWidth']):
                            #lock player played cards in
                            for card in tempLocal:
                                cardsPlayedLocal.add(card)

                    #if attacker is opponent
                    elif gameState == -1:
                        #potential defense, can be changed
                        if rectangle['LocalPlayer'] == True and inPlay(rectangle['TopLeftY'], received['Screen']['ScreenHeight'], received['Screen']['ScreenWidth']):
                            played.add(rectangle['CardCode'])
                
                #if the attacker has changed, the played cards are commited and added to final list
                if gameState != pastGameState and pastGameState != 0:
                    for card in tempLocal:
                        cardsPlayedLocal.add(card)
                    pastGameState = gameState
                    #print(cardsPlayedLocal)

                #commit these to a temp with intersect
                tempLocal = playedCardsChooser(tempLocal, played)
                

            elif received['GameState'] == 'Menus':
                gameState = yaml.safe_load(reqGame.text)
                #print(gameState)
                if gameState["GameID"] > gamesPlayed:
                    #print(gameState['LocalPlayerWon'])
                    #send completed game card set to server here
                    for card in tempLocal:
                        cardsPlayedLocal.add(card)

                    req = {
                        'win': gameState["LocalPlayerWon"],
                        'card_codes': list(cardsPlayedLocal),
                    }

                    jsonreq = json.dumps(req)

                    #if the first game of this session
                    if gamesPlayed == -1:
                        a = requests.post("http://localhost:8080/cards", data = jsonreq)
                        id = yaml.safe_load(a.text)["id"]
                        print("http://localhost:8080/view/"+ id)
                    else:
                        a = requests.put("http://localhost:8080/cards/" + id, data = jsonreq)
                        #print(a.status_code)
                    #print(cardsPlayedLocal)
                    gamesPlayed = gameState['GameID']
                    #reset sets for next game
                    cardsPlayedLocal = set()
        

main()