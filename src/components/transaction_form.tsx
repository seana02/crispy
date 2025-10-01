import { useForm, useStore } from "@tanstack/react-form";
import { Input } from "./Input";
import { useEffect, useRef, useState } from "react";
import { DayPicker, getDefaultClassNames } from "react-day-picker";
import 'react-day-picker/style.css';
import { PostingData } from "@/types";
import currency from "currency.js";
import { ConfirmationButton } from "./ConfirmationButton";
import { MdDelete, MdDeleteForever } from "react-icons/md";
import { Combobox } from "./Combobox";
import { useNavigate } from "@tanstack/react-router";
import { getTodayDatestring } from "@/routes/transactions_add";

interface TransactionFormProps {
    date: string,
    description: string,
    postings: {
        id: number,
        account: string,
        value: string,
        currency: string,
        comment: string,
    }[],
    onSubmit: (year: number, month: number, day: number, postings: PostingData[], desc: string, delete_list: number[]) => void,
    duplicated?: boolean
}

export function TransactionForm(props: TransactionFormProps) {
    const navigate = useNavigate();
    const transactionForm = useForm({
        defaultValues: {
            date: props.date,
            description: props.description,
            postings: props.postings,
        },
        onSubmit: (e: { value: { date: string, description: string, postings: PostingData[] } }) => {
            const datestring = new Date(e.value.date).toISOString();
            const year = +datestring.substring(0,4);
            const month = +datestring.substring(5,7);
            const day = +datestring.substring(8,10);
            let deleteList: number[] = [];
            for (let i in deleteObj) {
                if (deleteObj[i]) deleteList = [...deleteList, +i];
            }
            props.onSubmit(year, month, day, e.value.postings, e.value.description, deleteList);
        },
        validators: {
            onChange({ value }) {
                if (!value.date) return 'Enter a valid date';
                if (!value.description) return 'Enter a valid description';
                let sum = currency(0);
                if (value.postings.length <= 0) return 'Must have at least 1 posting';
                for (let a of value.postings) {
                    if (!a.account) return 'All Postings must have a valid account';
                    if (!a.value || isNaN(+a.value)) return 'All Postings must have a valid value';
                    if (!a.currency) return 'All Postings must specify a valid currency';
                    sum = sum.add(currency(a.value));
                }
                if (sum.intValue !== 0) return `Unbalanced transaction: ${sum.toString()}`;
                return undefined;
            }
        }
    });

    const formErrorMap = useStore(transactionForm.store, (state) => state.errorMap);

    const dialogRef = useRef<HTMLDialogElement>(null);

    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [month, setMonth] = useState(new Date(`${props.date}T00:00:00Z`));
    const [deleteObj, setDeleteObj] = useState<{[id: number]: boolean}>(
        function () {
            let deleteObj: {[id: number]: boolean} = {};
            for (let p of props.postings) {
                deleteObj[p.id] = false;
            }
            return deleteObj;
        }()
    );


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
                className={`bg-[rgba(0,0,0,0.3)] w-screen h-screen absolute top-0 left-0 z-10 ${isDialogOpen ? '' : 'hidden'}`}
                onClick={() => setIsDialogOpen(!isDialogOpen)}
            />
            <form className="flex flex-col max-h-full" onSubmit={e => {
                e.preventDefault();
                transactionForm.handleSubmit();
            }}>
                <div className="flex gap-4 mb-4">
                    <transactionForm.Field name="date">
                        {field => {
                            return (
                                <div className={`relative ${checkChange(props.duplicated ? getTodayDatestring() : props.date, field.state.value)}`}>
                                    <div
                                        // onClose={() => setIsDialogOpen(false)}
                                        className={`absolute top-[100%] bg-black border border-blue-400 rounded-lg p-2 z-20 ${isDialogOpen ? '' : 'hidden'}`}
                                    >
                                        <DayPicker
                                            month={month}
                                            onMonthChange={setMonth}
                                            autoFocus
                                            mode="single"
                                            timeZone="America/New_York"
                                            selected={new Date(`${field.state.value}T00:00:00-04:00`)}
                                            onSelect={d => field.handleChange(d.toISOString().substring(0, 10))}
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
                                        className="w-27 text-lg text-white"
                                        value={field.state.value}
                                        onChange={e => field.handleChange(e.target.value)}
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
                            <div className={`flex-1 flex ${checkChange(props.duplicated ? '' : props.description, field.state.value)}`}>
                                <Input
                                    type="text"
                                    id="description"
                                    name="description"
                                    placeholder="Description"
                                    className={`flex-1 ${checkValid(field.state.value === '')}`}
                                    value={field.state.value}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => field.handleChange(e.target.value)}
                                />
                            </div>
                        )}
                    </transactionForm.Field>
                    {
                        props.description === "" || props.duplicated ? <></> :
                        <button className="border rounded-md border-yellow-200 hover:border-yellow-400 px-1 text-center flex items-center justify-center" onClick={duplicateEntry}>
                            Duplicate
                        </button>
                    }
                </div>
                <transactionForm.Field name="postings" mode="array">
                    {field => (
                        <div className="flex-1 overflow-auto no-scrollbar">
                            <div className="flex gap-3 [&>*]:border-b-2 [&>*]:border-green-300 [&>*]:select-none">
                                <div className="w-8 text-right">#</div>
                                <div className="min-w-[249px] flex-1">Account</div>
                                <div className="min-w-[249px] text-right">Value</div>
                                <div className="min-w-0 w-12">CCY</div>
                                <div className="w-[30px] text-center">X</div>
                            </div>
                            {field.state.value.map((_, i) => (
                                <div key={i} className={`flex gap-3 my-3 ${props.duplicated ? 'bg-gray-700' : props.postings[i] ? '' : 'bg-gray-700'}`}>
                                    <div className="w-8 text-right select-none">{i+1}</div>
                                    <transactionForm.Field name={`postings[${i}].account`}>
                                        {subField => (
                                            <div className="min-w-0 flex-1">
                                                {/* <Input */}
                                                <Combobox
                                                    value={subField.state.value}
                                                    onChange={e => subField.handleChange(e.target.value)}
                                                    className={`w-full ${checkChange(props.postings[i]?.account, subField.state.value)} ${checkValid(subField.state.value === "")}`}
                                                    updateValue={str => subField.handleChange(str)}
                                                />
                                            </div>
                                        )}
                                    </transactionForm.Field>
                                    <transactionForm.Field
                                        name={`postings[${i}].value`}
                                    >
                                        {subField => (
                                            <div className="min-w-0">
                                                <Input
                                                    value={subField.state.value}
                                                    onChange={e => {
                                                        let input = e.target.value;
                                                        // cannot have anything except numbers, decimal, and dash
                                                        if (input.match(/[^0-9.-]/)) return;
                                                        // cannot have more than one decimal
                                                        if (input.match(/\.(?=.*\.)/)) return;
                                                        // toggle minus at front only
                                                        const minusCount = input.match(/-/g)?.length;
                                                        if (minusCount && minusCount % 2) {
                                                            input = "-" + input.replace(/-/g, "");
                                                        } else {
                                                            input = input.replace(/-/g, "");
                                                        }
                                                        subField.handleChange(input);
                                                    }}
                                                    className={`w-full text-right ${checkChange(props.postings[i]?.value, subField.state.value)} ${checkValid(subField.state.value === "" || isNaN(+subField.state.value))}`}
                                                />
                                            </div>
                                        )}
                                    </transactionForm.Field>
                                    <transactionForm.Field name={`postings[${i}].currency`} >
                                        {subField => (
                                            <div className="min-w-0 w-12">
                                                <Input
                                                    value={subField.state.value}
                                                    onChange={e => subField.handleChange(e.target.value)}
                                                    className={`w-full ${checkChange(props.postings[i]?.currency, subField.state.value)} ${checkValid(subField.state.value === "")}`}
                                                />
                                            </div>
                                        )}
                                    </transactionForm.Field>
                                    <ConfirmationButton
                                        className="border rounded-md border-red-200 hover:border-red-400 w-[30px] h-[30px] text-center flex items-center justify-center"
                                        onConfirm={() => {
                                            setDeleteObj({ ...deleteObj, [field.state.value[i].id]: true });
                                            field.removeValue(i);
                                        }}
                                        baseText={<MdDelete />}
                                        confirmationText={<MdDeleteForever />}
                                    />
                                </div>
                            ))}
                            <div
                                onClick={() => field.pushValue({ id: -1, account: '', value: '', currency: 'USD', comment: '' })}
                                className="border rounded-lg border-blue-200 hover:border-blue-400 w-full select-none text-center mt-3"
                            >
                                Add Posting
                            </div>
                        </div>
                    )}
                </transactionForm.Field>
                {formErrorMap.onChange ? (
                    <div className='border rounded-lg border-red-200 w-full my-4 select-none text-red-400 text-center'>
                        <em>Error: {formErrorMap.onChange}</em>
                    </div>
                ) : (
                    <button
                        onClick={transactionForm.handleSubmit}
                        type="button"
                        className="border rounded-lg border-red-200 hover:border-red-400 w-full my-4 select-none"
                    >
                        Submit
                    </button>
                    )}
            </form>
        </div>
    );

    function checkChange<T>(oldVal: T, newVal: T) {
        if (oldVal !== newVal) return 'bg-gray-700';
        if (!oldVal || !newVal) return '';
        return '';
    }

    function checkValid(check: boolean) {
        if (check) return 'border-red-200 hoover:border-red-400';
        return '';
    }

    function duplicateEntry() {
        console.log('hi');

        navigate({
            to: '../add',
            search: {
                date: props.date,
                description: props.description,
                postings: JSON.stringify(props.postings),
                duplicated: true
            }
        });
    }
}

