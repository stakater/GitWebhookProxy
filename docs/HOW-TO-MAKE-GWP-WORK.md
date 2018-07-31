# HOW TO MAKE GWP WORK

Below are the Steps to create a Git webhook, which can be used with GitwebhookProxy

1. Go to the `Settings` tab of your repository.
2. Click on `Webhooks` in left sidebar:
![alt text][webhooks]
3. Click on `Add webhook` button.
![alt text][add-webhook]
4. Enter `Payload URL`.
5. Enter `Content type` as `application/json`.
6. Enter your `Secret`.
![alt text][webhook-details]
7. Select the events to use with this webhook.
8. Check the `Active` checkbox.
9. Click on `Add webhook` button.
![alt text][save-webhook]
10. Use the `Secret` and `Payload URL` in GitWebHookProxy's Vanilla manifest.

NOTE: For further details regarding webhooks, please refer [here](https://developer.github.com/webhooks/).

[webhooks]: images/webhooks.png "Webhooks"
[add-webhook]: images/add-webhook.png "Add webhooks"
[webhook-details]: images/webhook-details.png "webhook-details"
[save-webhook]: images/save-webhook.png "Save webhook"