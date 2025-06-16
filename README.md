# lab-attendance
has the service and UI for the lab attendance touchpanels. The UI is in /web/ and shows prompts based off of websocket messages from the backend. The backend sends requests to the lab API and updates the UI by returning websockets with error and success messages. It also handles reads from cards from [hid-reader-microservice](https://github.com/byuoitav/hid-reader-microservice).

![image](https://github.com/user-attachments/assets/d0672bab-ec0a-4541-b2df-e9c56a149131)
