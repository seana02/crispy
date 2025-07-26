import { useQuery } from "@tanstack/react-query";
import { createRoute } from "@tanstack/react-router";
import { transactionRoute } from "./transactions";
import { PostingData, TransactionData } from "../types";
import { getTransactionDetails } from "../api";
import { useForm } from "@tanstack/react-form";
import { useEffect } from "react";

export const transactionDetailsRoute = createRoute({
    getParentRoute: () => transactionRoute,
    path: '/$id',
    component: TransactionDetails,
});

function TransactionDetails() {
    const transactionForm = useForm({
        defaultValues: {
            date: new Date(),
            description: "",
            postings: [] as PostingData[],
        },
        onSubmit: async ({ value }) => console.log(value)
    });

    const { id } = transactionDetailsRoute.useParams();
    const { isPending, isError, data, error } = useQuery<TransactionData>({
        queryKey: ['transaction_details', id],
        queryFn: () => getTransactionDetails({ id: +id }),
    });

    useEffect(() => {
        if (data) {
            transactionForm.reset({
                date: data.transaction_date,
                description: data.description,
                postings: data.postings,
            });
        }
    }, [data]);

    if (isPending) return <div>Pending...</div>
    if (isError) return <div>Error: {error.message}</div>

    return (
        <div className="flex flex-col">
            <form onSubmit={e => {
                e.preventDefault();
                transactionForm.handleSubmit();
            }}>
                <transactionForm.Field name="date">
                    {field => (
                        <div className="w-96">
                            <label htmlFor="date">Date: </label>
                            <input
                                type="text"
                                id="date"
                                name="date"
                                value={field.state.value.toString()}
                                onChange={e => field.handleChange(new Date(e.target.value))}
                            />
                        </div>
                    )}
                </transactionForm.Field>
                <transactionForm.Field name="description" >
                    {field => (
                        <div className="w-96">
                            <label htmlFor="description">Description: </label>
                            <input
                                type="text"
                                id="description"
                                name="description"
                                value={field.state.value}
                                onChange={e => field.handleChange(e.target.value)}
                            />
                        </div>
                    )}
                </transactionForm.Field>
                <transactionForm.Field name="postings" mode="array">
                    {field => (
                        <div>
                            {field.state.value.map((_, i) => (
                                <div key={i} className="flex flex-1 gap-2">
                                    <transactionForm.Field name={`postings[${i}].account`}>
                                        {subField => (
                                            <div className="flex gap-2 w-96">
                                                <label htmlFor={`postings[${i}].account`}>{'Account:'}</label>
                                                <input
                                                    value={subField.state.value}
                                                    onChange={e => subField.handleChange(e.target.value)}
                                                    className="flex-1"
                                                />
                                            </div>
                                        )}
                                    </transactionForm.Field>
                                    <transactionForm.Field name={`postings[${i}].value`}>
                                        {subField => (
                                            <div className="flex gap-2 w-24">
                                                <label htmlFor={`postings[${i}].account`}>{'Value:'}</label>
                                                <input
                                                    value={subField.state.value}
                                                    onChange={e => subField.handleChange(e.target.value)}
                                                    className="flex-1"
                                                />
                                            </div>
                                        )}
                                    </transactionForm.Field>
                                    <transactionForm.Field name={`postings[${i}].currency`}>
                                        {subField => (
                                            <div className="flex gap-2 w-24">
                                                <label htmlFor={`postings[${i}].account`}>{''}</label>
                                                <input
                                                    value={subField.state.value}
                                                    onChange={e => subField.handleChange(e.target.value)}
                                                    className="flex-1"
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

