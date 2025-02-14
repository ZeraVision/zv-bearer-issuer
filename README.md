# Bearer Issuer

A simple API application designed to interact with the ZV Indexer API for issuing Bearer tokens linked to your API key.

## How It Works

To configure the application, provide the following environment variables:

```plaintext
API_KEY=
API_SECRET=
API_PORT=
INDEXER_URL=
```

- Modify the authentication logic as needed.
- Adjust values or implement logic to determine limits for bearer tokens.

### Default Limits:

- **Maximum Requests Per Second (RPS):** 1
- **Maximum Burst Per Second (BPS):** 5
- **Maximum Requests Per Day (RPD):** 1000

Both Bearer tokens and API keys operate using the **Token Bucket Algorithm** for rate limiting.

## Obtaining an API Key

To request an API key, please contact us at [Zera Vision](https://www.zera.vision/contact).

A developer platform with automated API key issuance and native GUI analytics is planned for future release. Stay updated at [ZV Explorer](https://explorer.zera.vision/apis) *(coming soon)*.

## Api Usage
To request get a bearer from this sample API, you simply need to issue the below request. Modify as needed.
```
curl --location --request POST '{{your_bearer_endpoint}}/store?requestType=getBearer'
```

Sample response:
```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJicHMiOjUsImV4cCI6MTczOTY1NjA3MSwianRpIjoiZWYwMDRlNWQtNmVlZS00MWI4LTlkOGItYzA3ZTUzYjcxMzhmIiwicnBkIjoxMDAwLCJycHMiOjF9.gNVu85zWBx4xGOA4TEUwnzXPc3EqiO2k-X4BGWS-3Pc",
    "validUntil": 1739656071,
    "maxRequestsPerSecond": 1,
    "maxBurstPerSecond": 5,
    "maxRequestsPerDay": 1000
}
```

## Bearer Usage on Zv-Indexer
In the `Authorization` header of a request to the zv-indexer include `Bearer <token>` format. For example:
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJicHMiOjUsImV4cCI6MTczOTY1NjA3MSwianRpIjoiZWYwMDRlNWQtNmVlZS00MWI4LTlkOGItYzA3ZTUzYjcxMzhmIiwicnBkIjoxMDAwLCJycHMiOjF9.gNVu85zWBx4xGOA4TEUwnzXPc3EqiO2k-X4BGWS-3Pc
```

This will allow the user to make requests to the indexer as long as the related authorization is valid and within rate limits. If the bearer expires, is unauthorized, or the API Key is no longer valid, status code 401 will be returned. In the case of a simple expiry, your application may wish to retrieve a new bearer token from your issuer.
