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
  lastStatusCode: number;
  createdAt: string;
  updatedAt: string;
};

export type GetWebhooks = {
  webhooks: {
    nodes: [Webhook];
    totalCount: number;
    pageInfo: {
      startCursor?: string;
      endCursor?: string;
      hasPreviousPage: boolean;
      hasNextPage: boolean;
    };
  };
};

export const GET_WEBHOOKS = gql`
  query GetWebhooks($first: Uint!) {
    webhooks(first: $first, filter: {}) {
      nodes {
        id
        url
        method
        lastStatusCode
        createdAt
      }
      totalCount
      pageInfo {
        startCursor
        endCursor
        hasPreviousPage
        hasNextPage
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
      lastStatusCode
      attempts {
        body
        headers {
          key
          value
        }
        elapsedTime
        createdAt
      }
      createdAt
      updatedAt
    }
  }
`;
