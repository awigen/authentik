import { MessageLevel } from "@goauthentik/common/messages";
import { showMessage } from "@goauthentik/elements/messages/MessageContainer";

import { t } from "@lingui/macro";

import { Middleware, ResponseContext } from "@goauthentik/api";

export class MessageMiddleware implements Middleware {
    post(context: ResponseContext): Promise<Response | void> {
        if (context.response.status >= 500) {
            showMessage({
                level: MessageLevel.error,
                message: t`API request failed`,
                description: `${context.init.method} ${context.url}: ${context.response.status}`,
            });
        }
        return Promise.resolve(context.response);
    }
}
