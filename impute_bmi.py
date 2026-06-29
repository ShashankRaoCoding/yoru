#!/usr/bin/env python3
"""
Impute missing y values using a linear regression of known (x, y) pairs.

For entities with multiple rows (same uid), the known x/y values are pooled
per uid before fitting, so each unique entity contributes equally regardless
of how many duplicate rows it has.

Usage:
    python impute_y.py <input_csv> <uid_col> <x_col> <y_col> [--output <path>]

Example:
    python impute_y.py data.csv entity_id revenue headcount --output data_imputed.csv
"""

import argparse
import sys
import io
import pandas as pd
import numpy as np
from sklearn.linear_model import LinearRegression


def parse_args():
    parser = argparse.ArgumentParser(
        description="Impute missing y values from x values, respecting uid grouping."
    )
    parser.add_argument("input", help="Path to the input CSV file")
    parser.add_argument("uid_col", help="Column name for the unique entity identifier")
    parser.add_argument("x_col", help="Column name for the predictor (x)")
    parser.add_argument("y_col", help="Column name for the target to impute (y)")
    parser.add_argument(
        "--output", "-o",
        default=None,
        help="Path for the output CSV (default: overwrites input with '_imputed' suffix)"
    )
    parser.add_argument(
        "--method", "-m",
        choices=["linear", "median", "mean"],
        default="linear",
        help=(
            "Imputation method: "
            "'linear' fits a linear regression (default), "
            "'median' uses the median y/x ratio, "
            "'mean' uses the mean y value directly (ignores x)"
        )
    )
    return parser.parse_args()


def load_and_validate(path, uid_col, x_col, y_col):
    if path == "":
        df = pd.read_csv(io.StringIO(sys.stdin.read()))
    else:
        df = pd.read_csv(path)  # was: bare `import pandas` instead of reading the file

    missing = [c for c in [uid_col, x_col, y_col] if c not in df.columns]
    if missing:
        sys.exit(f"ERROR: Column(s) not found in CSV: {missing}\n"
                 f"Available columns: {df.columns.tolist()}")

    return df


def get_uid_level_data(df, uid_col, x_col, y_col):
    """
    Collapse multiple rows per uid into one representative row.
    For each uid, take the first non-null x and first non-null y value.
    A uid 'has known y' if any of its rows has a non-null y.
    """
    def first_valid(s):
        valid = s.dropna()
        return valid.iloc[0] if len(valid) > 0 else np.nan

    uid_df = (
        df.groupby(uid_col, sort=False)[[x_col, y_col]]
        .agg(first_valid)
        .reset_index()
    )
    return uid_df


def fit_and_predict(uid_df, x_col, y_col, method):
    known = uid_df.dropna(subset=[y_col])
    unknown = uid_df[uid_df[y_col].isna() & uid_df[x_col].notna()]

    print(f"\nUID summary:")
    print(f"  UIDs with known y    : {len(known)}")
    print(f"  UIDs needing imputation (y missing, x present): {len(unknown)}")
    print(f"  UIDs missing both x and y: "
          f"{uid_df[uid_df[y_col].isna() & uid_df[x_col].isna()].shape[0]}")

    if len(known) == 0:
        sys.exit("ERROR: No rows have a known y value — cannot impute.")

    if len(unknown) == 0:
        print("\nNothing to impute — all rows with x already have y.")
        return uid_df

    if method == "linear":
        if len(known) < 2:
            print("WARNING: Only 1 known y value — falling back to mean imputation.")
            method = "mean"
        else:
            X_train = known[[x_col]].values
            y_train = known[y_col].values
            model = LinearRegression()
            model.fit(X_train, y_train)
            r2 = model.score(X_train, y_train)
            print(f"\n  Linear regression fit: R² = {r2:.4f}, "
                  f"slope = {model.coef_[0]:.4f}, "
                  f"intercept = {model.intercept_:.4f}")

            X_pred = unknown[[x_col]].values
            predictions = model.predict(X_pred)
            uid_df.loc[unknown.index, y_col] = predictions

    if method == "median":
        known_ratio = known[y_col] / known[x_col].replace(0, np.nan)
        ratio = known_ratio.median()
        print(f"\n  Median y/x ratio used: {ratio:.4f}")
        uid_df.loc[unknown.index, y_col] = unknown[x_col] * ratio

    if method == "mean":
        mean_y = known[y_col].mean()
        print(f"\n  Mean y used for imputation: {mean_y:.4f}")
        uid_df.loc[unknown.index, y_col] = mean_y

    return uid_df


def propagate_back(df, uid_df, uid_col, y_col):
    """
    Map imputed y values back to every row in the original dataframe.
    Only fills rows where y was originally missing.
    """
    imputed_map = uid_df.set_index(uid_col)[y_col]

    # Track which rows were originally missing
    originally_missing = df[y_col].isna()

    # Fill missing y using the imputed uid-level value
    df[y_col] = df.apply(
        lambda row: imputed_map.get(row[uid_col], row[y_col])
        if pd.isna(row[y_col])          # was: `import pandas` (syntax error)
        else row[y_col],
        axis=1
    )

    n_filled = originally_missing.sum() - df[y_col].isna().sum()
    print(f"\n  Rows filled: {n_filled} "   # was: bare `import pandas` (syntax error)
          f"(out of {originally_missing.sum()} originally missing)")

    return df


def main():
    args = parse_args()

    print(f"Loading: {args.input}")
    df = load_and_validate(args.input, args.uid_col, args.x_col, args.y_col)
    print(f"Shape: {df.shape}")
    print(f"Missing y values: {df[args.y_col].isna().sum()}")

    # Collapse to uid-level for fitting
    uid_df = get_uid_level_data(df, args.uid_col, args.x_col, args.y_col)

    # Fit model and predict missing uid-level y
    uid_df = fit_and_predict(uid_df, args.x_col, args.y_col, args.method)

    # Propagate imputed values back to all rows
    df = propagate_back(df, uid_df, args.uid_col, args.y_col)

    # Write output
    if args.output is None:
        base = args.input.rsplit(".", 1)
        args.output = f"{base[0]}_imputed.{base[1]}" if len(base) == 2 else args.input + "_imputed"

    df.to_csv(args.output, index=False)   # was: print(df.to_csv(...)) which printed to stdout instead of saving
    print(f"\nSaved to: {args.output}")
    remaining_missing = df[args.y_col].isna().sum()
    if remaining_missing:
        print(f"Note: {remaining_missing} rows still have missing y "
              f"(their uid also had no x value to predict from).")


if __name__ == "__main__":
    main()
