import { useState } from "react";

interface ConfirmationProps {
    onConfirm: () => void,
    className: string,
    baseText: string | React.ReactNode,
    confirmationText: string | React.ReactNode
}

export function ConfirmationButton({ onConfirm, className, baseText, confirmationText }: ConfirmationProps) {
    let [ clicked, setClicked ] = useState(false);
    let [ buffer, setBuffer ] = useState(false);

    return (clicked && buffer) ? 
    (
        <div
            className={className}
                onClick={() => {
                    setClicked(false);
                    setBuffer(false);
                    onConfirm();
                }}
        >
            {confirmationText}
        </div>
    ) : (
        <div
            className={className}
            onClick={() => {
                setClicked(true);
                setTimeout(() => setBuffer(true), 10);
                setTimeout(() => {
                    setClicked(false);
                    setBuffer(false);
                }, 2000);
            }}
        >
            {baseText}
        </div>
    );
}
