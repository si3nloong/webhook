import { gql } from "@apollo/client";

export const GET_WEBHOOKS = gql`
    query GetWebhooks {
        webhooks {
            nodes {
                id
                url
                method
            }
        }
    }
`;