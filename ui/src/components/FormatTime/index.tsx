/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import { FC, memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

import classNames from 'classnames';
import dayjs from 'dayjs';

import { pathFactory } from '@/router/pathFactory';

interface Props {
  time: number;
  className?: string;
  preFix?: string;
  questionId?: string | null;
  slugTitle?: string;
  answerId?: string | null;
  operation_type?: string;
}

const Index: FC<Props> = ({
  questionId,
  slugTitle,
  answerId,
  time,
  preFix,
  className,
  operation_type,
}) => {
  const { t } = useTranslation();
  const formatTime = (from) => {
    const now = Math.floor(dayjs().valueOf() / 1000);
    const between = now > from ? now - from : 0;

    if (between <= 1) {
      return t('dates.now');
    }
    if (between > 1 && between < 60) {
      return t('dates.x_seconds_ago', { count: between });
    }

    if (between >= 60 && between < 3600) {
      const min = Math.floor(between / 60);
      return t('dates.x_minutes_ago', { count: min });
    }
    if (between >= 3600 && between < 3600 * 24) {
      const h = Math.floor(between / 3600);
      return t('dates.x_hours_ago', { count: h });
    }

    if (
      between >= 3600 * 24 &&
      between < 3600 * 24 * 366 &&
      dayjs.unix(from).format('YYYY') === dayjs.unix(now).format('YYYY')
    ) {
      return dayjs.unix(from).tz().format(t('dates.long_date'));
    }

    return dayjs.unix(from).tz().format(t('dates.long_date_with_year'));
  };

  if (!time) {
    return null;
  }

  // return (
  //   // <NavLink
  //   //                   to={pathFactory.answerLanding({
  //   //                     questionId: li.id,
  //   //                     slugTitle: li.url_title,
  //   //                     answerId: li.last_answer_id,
  //   //                   })}
  //   //                   className="link-dark">
  //   //                   {/* {li.title}
  //   //                   {li.status === 2 ? ` [${t('closed')}]` : ''} */}
  //   //                 </NavLink>
  //   <time
  //     className={classNames('', className)}
  //     dateTime={dayjs.unix(time).tz().toISOString()}
  //     title={dayjs.unix(time).tz().format(t('dates.long_date_with_time'))}>
  //     {preFix ? `${preFix} ` : ''}
  //     {formatTime(time)}
  //   </time>
  // );

  const timeElement = (
    <time
      className={classNames('', className)}
      dateTime={dayjs.unix(time).tz().toISOString()}
      title={dayjs.unix(time).tz().format(t('dates.long_date_with_time'))}>
      {preFix ? `${preFix} ` : ''}
      {formatTime(time)}
    </time>
  );

  if (operation_type === 'answered') {
    return (
      <Link
        to={pathFactory.answerLanding({
          questionId: questionId || '',
          slugTitle: slugTitle || '',
          answerId: answerId || '',
        })}>
        {timeElement}
      </Link>
    );
  }

  return timeElement;
};

export default memo(Index);
