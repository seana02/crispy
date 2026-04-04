import { ParentProps, Show, children } from "solid-js";
import { Portal } from "solid-js/web";

import "../styles/globals.css";

interface ConfirmationPopupProps {
    isOpen: boolean;
    onClose: () => void;
}

export default function ConfirmationPopup(props: ParentProps<ConfirmationPopupProps>) {
    const safeChildren = children(() => props.children);

    return (
        <Show when={props.isOpen}>
            <Portal>
                <div
                    class="overlay"
                    onClick={props.onClose}
                >
                    <div
                        onClick={(e) => e.stopPropagation()}
                        class="popup"
                    >
                        {safeChildren()}
                    </div>
                </div>
            </Portal>
        </Show>
    );
}
