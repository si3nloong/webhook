import { gql } from "@apollo/client";

type HttpHeader = [{ key: string; value: string }];

type WebhookAttempt = {
  body: string;
  headers: HttpHeader;
  elapsedTime: number;
  createdAt: string;
};

export type Webhook = {
  id: string;
  method: string;
  url: string;
  body: string;
  headers: HttpHeader;
  timeout: number;
  noOfRetries: number;
  attempts: [WebhookAttempt];
  latestStatusCode: number;
  createdAt: string;
  updatedAt: string;
};

export const GET_WEBHOOKS = gql`
  query GetWebhooks {
    webhooks {
      nodes {
        id
        url
        method
        latestStatusCode
        createdAt
      }
    }
  }
`;

export const FIND_WEBHOOK = gql`
  query FindWebhookByID($id: ID!) {
    webhook(id: $id) {
      id
      url
      method
      headers {
        key
        value
      }
      body
      timeout
      attempts {
        body
        elapsedTime
        createdAt
      }
      createdAt
      updatedAt
    }
  }
`;
