import { createEffect, createMemo, createUniqueId, onCleanup, Show } from 'solid-js';
import './ConfirmDialog.css';

interface ConfirmDialogProps {
    isOpen: boolean;
    title: string;
    message: string;
    confirmText?: string;
    cancelText?: string;
    variant?: 'danger' | 'warning' | 'default';
    onConfirm: () => void;
    onCancel: () => void;
    error?: string | null;
    isLoading?: boolean;
}

export function ConfirmDialog(props: ConfirmDialogProps) {
    const titleId = createUniqueId();

    createEffect(() => {
        if (!props.isOpen) return;
        const onKey = (e: KeyboardEvent) => { if (e.key === 'Escape') props.onCancel(); };
        document.addEventListener('keydown', onKey);
        onCleanup(() => document.removeEventListener('keydown', onKey));
    });

    const confirmClass = createMemo(() => {
        switch (props.variant) {
            case 'danger': return 'btn btn-danger';
            case 'warning': return 'btn btn-primary';
            default: return 'btn btn-primary';
        }
    });

    const icon = () => props.variant === 'default' ? 'ℹ️' : '⚠️';

    return (
        <Show when={props.isOpen}>
            <div class="confirm-overlay">
                <div class="confirm-backdrop" onClick={props.onCancel} />
                <div
                    class="confirm-dialog"
                    role="dialog"
                    aria-modal="true"
                    aria-labelledby={titleId}
                    tabindex="-1"
                    ref={(el) => queueMicrotask(() => el.focus())}
                >
                    <div class="confirm-header">
                        <span class="confirm-icon">{icon()}</span>
                        <h3>{props.title}</h3>
                    </div>
                    <div class="confirm-body">
                        <Show when={props.error}>
                            <p role="alert" class="confirm-error">{props.error}</p>
                        </Show>
                        <p>{props.message}</p>
                    </div>
                    <div class="confirm-footer">
                        <button class="btn btn-secondary" onClick={props.onCancel} type="button">
                            {props.cancelText ?? 'Cancel'}
                        </button>
                        <button class={confirmClass()} onClick={props.onConfirm} type="button">
                            {props.confirmText ?? 'Confirm'}
                        </button>
                    </div>
                </div>
            </div>
        </Show>
    );
}
