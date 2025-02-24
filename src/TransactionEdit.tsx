import { RefObject, useEffect, useRef, useState } from "react";
import DatePicker from "react-datepicker";
import 'react-datepicker/dist/react-datepicker.css';
import './styles/Transaction.css';

interface PostingData {
    acct: string;
    val: number;
    currency: string;
    comment: string;
}

interface TransactionData {
    date: Date;
    postings: PostingData[];
    desc: string;
}

interface TEditProps {
    data?: TransactionData;
    updateTab: (newTab: string) => void;
}

export default function TransactionEdit(props: TEditProps) {

    const [data, setData] = useState<TransactionData> ({
        date: props.data?.date || new Date(),
        postings: props.data?.postings || [],
        desc: props.data?.desc || ''
    });

    const endRef: RefObject<HTMLDivElement> = useRef(null);

    const change = (e: React.ChangeEvent<HTMLInputElement>) => {
        const {name, value} = e.target;
        setData(prevData => ({...prevData, [name]: value}));
    }

    const submit = (e: React.ChangeEvent<HTMLFormElement>) => {
        e.preventDefault();
        console.log(data);
    };

    const changePosting = (e: React.ChangeEvent<HTMLInputElement>, id: number) => {
        const { name, value } = e.target;
        if (name.split("-")[1] === "val") {
            const regex = /-?([1-9][0-9]*\.?[0-9]*)?/;
            const m = value.match(regex);
            if (value.length > 0 && (!m || m[0].length !== value.length)) return;
        }
        let newPostings = data.postings.map((v, i) => {
            if (i !== id) return v;
            return {
                ...v,
                [name.split("-")[1]]: value
            }
        });
        setData(prevData => ({...prevData, postings: newPostings}));
    }

    let postings = [];
    for (let i = 0; i < data.postings.length; i++) {
        const ch = (e: React.ChangeEvent<HTMLInputElement>) => changePosting(e, i);
        postings.push(
            <PostingRow 
                i={i}
                endRef={endRef}
                postings={data.postings}
                ch={ch}
                delete={() => setData(prevData => ({...prevData, postings: [...prevData.postings.filter((_, id) => id !== i)]}))}
            />
        );
    }
    
    return (
        <div id="transaction-edit">
            <h1>Edit Transaction</h1>
            <form id="transaction-edit-form" onSubmit={submit}>

                <div className="form-row">
                    <div className="form-element">
                        <label htmlFor="date">Date:</label>
                        <DatePicker name="date" selected={data.date} onChange={(date: Date | null) => setData(prevData => ({...prevData, date: date || new Date()}))} />
                    </div>

                    <div className="form-element flex1">
                        <label htmlFor="desc">Description:</label>
                        <input type="text" name="desc" id="form-desc" onChange={change} value={data.desc} />
                    </div>
                </div>

                <div id="form-button-row" className="form-row">
                    <button id="form-add-posting" onClick={() => {
                        setData(prevData => ({
                            ...prevData,
                            postings: [
                                ...prevData.postings,
                                {
                                    acct: "",
                                    val: 0,
                                    currency: "",
                                    comment: ""
                                }
                            ]
                        }));
                        if (endRef.current) {
                            endRef.current.scrollIntoView({ behavior: "smooth" });
                        }
                    }}>
                        New Posting
                    </button>

                    <input type="submit" id="form-submit" formAction="submit" value="Create" />
                </div>

                <div id="form-postings-list">
                    {postings}
                </div>

            </form>
        </div>
    );

}

interface PostingRowProps {
    i: number;
    endRef: RefObject<HTMLDivElement>;
    postings: PostingData[];
    ch: (e: React.ChangeEvent<HTMLInputElement>) => void;
    delete: () => void;
}

function PostingRow(props: PostingRowProps) {
    const lbl = `posting${props.i}`;
    return (
        <div key={props.i} ref={props.i+1 === props.postings.length ? props.endRef : null} className="form-row">
            <div className="form-element" style={{ width: "30px" }}>{props.i+1}.</div>

            <div className="form-element flex4">
                <label htmlFor={`${lbl}-acct`}>Account:</label>
                <input type="text" name={`${lbl}-acct`} className="form-acct" onChange={props.ch} value={props.postings[props.i]?.acct || ""}/>
            </div>

            <div className="form-element flex2">
                <label htmlFor={`${lbl}-val`}>Value:</label>
                <input type="text" name={`${lbl}-val`} className="form-val" onChange={props.ch} value={props.postings[props.i]?.val || ""}/>
            </div>

            <div className="form-element flex1">
                <label htmlFor={`${lbl}-currency`}>Currency:</label>
                <input type="text" name={`${lbl}-currency`} className="form-currency" onChange={props.ch} value={props.postings[props.i]?.currency || ""}/>
            </div>

            <div className="form-element flex4">
                <label htmlFor={`${lbl}-comment`}>Comment:</label>
                <input type="text" name={`${lbl}-comment`} className="form-comment" onChange={props.ch} value={props.postings[props.i]?.comment || ""}/>
            </div>

            <div className="form-element form-delete" onClick={props.delete}>
                <div>X</div>
            </div>
        </div>
    );
}
