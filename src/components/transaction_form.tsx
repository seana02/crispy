import { useForm } from "@tanstack/react-form";
import { Input } from "./Input";

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
    return (
        <div className="flex flex-col h-full overflow-y-scroll p-10 text-xl">
            <form onSubmit={e => {
                e.preventDefault();
                transactionForm.handleSubmit();
            }}>
                <div className="flex gap-4 mb-4">
                    <transactionForm.Field name="date">
                        {field => (
                            <Input
                                type="date"
                                id="date"
                                name="date"
                                className="w-30 text-lg"
                                value={field.state.value instanceof Date ? field.state.value.toISOString().slice(0, 10) : ''}
                                onChange={e => { console.log(e.target.value); field.handleChange(new Date(e.target.value)) }}
                            />
                        )}
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
                                <div className="min-w-[249px] grow-[5]">Account</div>
                                <div className="min-w-[249px] grow-[1]">Value</div>
                                <div className="min-w-0 w-12">CCY</div>
                            </div>
                            {field.state.value.map((_, i) => (
                                <div key={i} className="flex gap-4 my-3">
                                    <div className="w-8 text-right">{i}</div>
                                    <transactionForm.Field name={`postings[${i}].account`}>
                                        {subField => (
                                            <div className="min-w-0 grow-[5]">
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
                                </div>
                            ))}
                        </div>
                    )}
                </transactionForm.Field>
            </form>
        </div>
    );
}
