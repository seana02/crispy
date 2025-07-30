import { useForm } from "@tanstack/react-form";
import { Input } from "./Input";
import { useEffect, useRef, useState } from "react";
import { DayPicker, getDefaultClassNames } from "react-day-picker";
import 'react-day-picker/style.css';

interface TransactionFormProps {
    date: Date,
    description: string,
    postings: {
        account: string,
        value: string,
        currency: string,
        comment: string,
    }[],
    onSubmit: () => void,
}

export function TransactionForm(props: TransactionFormProps) {
    const transactionForm = useForm({
        defaultValues: {
            date: props.date,
            description: props.description,
            postings: props.postings,
        },
        onSubmit: props.onSubmit,
    });

    const dialogRef = useRef<HTMLDialogElement>(null);

    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [month, setMonth] = useState(props.date);

    useEffect(() => {
        if (!dialogRef.current) return;
        if (isDialogOpen) {
            dialogRef.current.showModal();
        } else {
            dialogRef.current.close();
        }
    }, [isDialogOpen]);

    return (
        <div className="flex flex-col h-full overflow-y-scroll p-10 text-xl">
            <div
                className={`bg-[rgba(0,0,0,0.3)] w-screen h-screen absolute top-0 left-0 ${isDialogOpen ? '' : 'hidden'}`}
                onClick={() => setIsDialogOpen(!isDialogOpen)}
            />
            <form onSubmit={e => {
                e.preventDefault();
                transactionForm.handleSubmit();
            }}>
                <div className="flex gap-4 mb-4">
                    <transactionForm.Field name="date">
                        {field => {
                            return (
                                <div className="relative">
                                    <div
                                        // onClose={() => setIsDialogOpen(false)}
                                        className={`absolute top-[100%] bg-black border border-blue-400 rounded-lg p-2 ${isDialogOpen ? '' : 'hidden'}`}
                                    >
                                        <DayPicker
                                            month={month}
                                            onMonthChange={setMonth}
                                            autoFocus
                                            mode="single"
                                            timeZone="UTC"
                                            selected={field.state.value}
                                            onSelect={d => field.handleChange(d as Date)}
                                            classNames={{
                                                today: 'text-green-200',
                                                selected: 'rdp-selected [&>*]:border-blue-200!',
                                                chevron: 'fill-blue-400',
                                                month_caption: `${getDefaultClassNames().month_caption} mx-3`
                                            }}
                                            required
                                        />
                                    </div>
                                    <Input
                                        type="date"
                                        id="date"
                                        name="date"
                                        className="w-30 text-lg text-white"
                                        value={new Date(field.state.value).toISOString().slice(0, 10)}
                                        onChange={e => field.handleChange(new Date(e.target.value))}
                                        onClick={e => {
                                            e.preventDefault();
                                            setIsDialogOpen(!isDialogOpen);
                                        }}
                                    />
                                </div>
                            );
                        }}
                    </transactionForm.Field>
                    <transactionForm.Field name="description" >
                        {field => (
                            <div className="flex-1 flex">
                                <Input
                                    type="text"
                                    id="description"
                                    name="description"
                                    placeholder="Description"
                                    className="flex-1"
                                    value={field.state.value}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => field.handleChange(e.target.value)}
                                />
                            </div>
                        )}
                    </transactionForm.Field>
                </div>
                <transactionForm.Field name="postings" mode="array">
                    {field => (
                        <div>
                            <div className="flex gap-4 [&>*]:border-b-2 [&>*]:border-green-300">
                                <div className="w-8 text-right">ID</div>
                                <div className="min-w-[249px] grow-[8]">Account</div>
                                <div className="min-w-[249px] grow-[1] text-right">Value</div>
                                <div className="min-w-0 w-12">CCY</div>
                                <div className="w-[30px] mx-[-8px] text-center">X</div>
                            </div>
                            {field.state.value.map((_, i) => (
                                <div key={i} className="flex gap-4 my-3">
                                    <div className="w-8 text-right">{i}</div>
                                    <transactionForm.Field name={`postings[${i}].account`}>
                                        {subField => (
                                            <div className="min-w-0 grow-[8]">
                                                <Input
                                                    value={subField.state.value}
                                                    onChange={e => subField.handleChange(e.target.value)}
                                                    className="w-full"
                                                />
                                            </div>
                                        )}
                                    </transactionForm.Field>
                                    <transactionForm.Field name={`postings[${i}].value`}>
                                        {subField => (
                                            <div className="min-w-0 grow-[1]">
                                                <Input
                                                    value={subField.state.value}
                                                    onChange={e => subField.handleChange(e.target.value)}
                                                    className="w-full text-right"
                                                />
                                            </div>
                                        )}
                                    </transactionForm.Field>
                                    <transactionForm.Field name={`postings[${i}].currency`}>
                                        {subField => (
                                            <div className="min-w-0 w-12">
                                                <Input
                                                    value={subField.state.value}
                                                    onChange={e => subField.handleChange(e.target.value)}
                                                    className="w-full"
                                                />
                                            </div>
                                        )}
                                    </transactionForm.Field>
                                    <ConfirmationButton
                                        className="border rounded-md border-red-200 hover:border-red-400 w-[30px] h-[30px] mx-[-8px]"
                                        onConfirm={() => field.removeValue(i)}
                                        baseText={"?"}
                                        confirmationText={"X"}
                                    />
                                </div>
                            ))}
                            <button
                                onClick={() => field.pushValue({ account: '', value: '', currency: 'USD', comment: '' })}
                                type="button"
                                className="border rounded-lg border-blue-200 hover:border-blue-400 w-full"
                            >
                                Add Posting
                            </button>
                        </div>
                    )}
                </transactionForm.Field>
            </form>
        </div>
    );
}

function ConfirmationButton({ onConfirm, className, baseText, confirmationText }: { onConfirm: () => void, className: string, baseText: string, confirmationText: string }) {
    let [ clicked, setClicked ] = useState(false);
    let [ buffer, setBuffer ] = useState(false);

    return (clicked && buffer) ? 
    (
        <button
            className={className}
                onClick={() => {
                    setClicked(false);
                    setBuffer(false);
                    onConfirm();
                }}
        >
            {confirmationText}
        </button>
    ) : (
        <button
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
        </button>
    );
}
