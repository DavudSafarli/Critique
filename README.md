# Critique

Open-source tool ( call it a microservice :wink: ) for managing customer feedbacks.
Critique **will be** able up and running in few minutes ( I hope ) , using the docker image.

API will have the following endpoints. Updating exising feedbacks and tags will also be available.

| Method | Endpoint | Description |
| ------ | :------- | :---------- |
| GET    | /v1/feedbacks     | get with pagination                      |
| GET    | /v1/feedbacks/:id | get feedback details                     |
| POST   | /v1/feedbacks     | post a feedback                          |
| GET    | /v1/tags          | get all tags(complaint, proposal, etc)   |
| POST   | /v1/tags          | create tags                              |
| PUT    | /v1/tags/merge    | merge tag with another                   |

## Database Structure

<table>
<tr>
<td valign="top">
    <table>
        <tr>
            <th colspan="2">feedbacks</th>
        </tr>
        <tr>
            <td><i>column name</u></td>
            <td><i>type</u></td>
        </tr>
        <tr>
            <td><b>id</b></td>
            <td>uint</td>
        </tr>
        <tr>
            <td><b>title</b></td>
            <td>string</td>
        </tr>
        <tr>
            <td><b>body</b></td>
            <td>string</td>
        </tr>
        <tr>
            <td><b>created_by</b></td>
            <td>string</td>
        </tr>
    </table>
</td>


<td valign="top">
    <table>
        <tr>
            <th colspan="2">attachments</th>
        </tr>
        <tr>
            <td><i>column name</u></td>
            <td><i>type</u></td>
        </tr>
        <tr>
            <td><b>id</b></td>
            <td>uint</td>
        </tr>
        <tr>
            <td><b>name</b></td>
            <td>string</td>
        </tr>
        <tr>
            <td><b>path</b></td>
            <td>string</td>
        </tr>
        <tr>
            <td><b>feedback_id</b></td>
            <td>uint</td>
        </tr>
    </table>
</td>


<td valign="top">
    <table>
        <tr>
            <th colspan="2">tags</th>
        </tr>
        <tr>
            <td><i>column name</u></td>
            <td><i>type</u></td>
        </tr>
        <tr>
            <td><b>id</b></td>
            <td>uint</td>
        </tr>
        <tr>
            <td><b>name</b></td>
            <td>string</td>
        </tr>
    </table>
</td>


<td valign="top">
    <table>
        <tr>
            <th colspan="2">feedbacks_types</th>
        </tr>
        <tr>
            <td><i>column name</u></td>
            <td><i>type</u></td>
        </tr>
        <tr>
            <td><b>feedback_id</b></td>
            <td>uint</td>
        </tr>
        <tr>
            <td><b>tag_id</b></td>
            <td>uint</td>
        </tr>
    </table>
</td>
</tr>
</table>


